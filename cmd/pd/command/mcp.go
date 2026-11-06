// mcp.go
//
// Copyright (c) 2018-2026 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package command

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"filippo.io/age"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/urfave/cli/v2"

	"github.com/jkawamoto/go-pixeldrain"
	"github.com/jkawamoto/go-pixeldrain/client/file"
	"github.com/jkawamoto/go-pixeldrain/client/list"
	"github.com/jkawamoto/go-pixeldrain/cmd/pd/auth"
	"github.com/jkawamoto/go-pixeldrain/cmd/pd/status"
	"github.com/jkawamoto/go-pixeldrain/models"
)

type mcpTextInput struct {
	Text       string   `json:"text" jsonschema:"Text to upload"`
	Name       string   `json:"name,omitempty" jsonschema:"File name (defaults to text.txt)"`
	Recipients []string `json:"recipients,omitempty" jsonschema:"Recipient IDs from list_recipients; omit for an unencrypted upload"`
}

type mcpFile struct {
	Name string `json:"name" jsonschema:"File name"`
	Data string `json:"data" jsonschema:"Base64 encoded file data, data URI, or local file path"`
}

type mcpFilesInput struct {
	Files      []mcpFile `json:"files" jsonschema:"Files to upload"`
	Recipients []string  `json:"recipients,omitempty" jsonschema:"Recipient IDs from list_recipients; omit for an unencrypted upload"`
}

type mcpURLInput struct {
	URL  string `json:"url" jsonschema:"Pixeldrain file download URL"`
	Path string `json:"path,omitempty" jsonschema:"Local output path; defaults to the downloaded file name"`
}

type mcpRecipient struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type mcpRecipientsResult struct {
	Recipients []mcpRecipient `json:"recipients"`
}

type mcpURLResult struct {
	URL string `json:"url"`
}

type mcpPathResult struct {
	Path string `json:"path"`
}

func CmdMCP(c *cli.Context) error {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	server, err := newMCPServer(c)
	if err != nil {
		return cli.Exit(err, status.InvalidArgument)
	}
	server.AddReceivingMiddleware(logMCPCalls(logger))
	logger.Info("MCP server started")
	err = server.Run(context.Background(), &mcp.StdioTransport{})
	if err != nil {
		logger.Error("MCP server stopped with an error")
	} else {
		logger.Info("MCP server stopped")
	}
	return err
}

func newMCPServer(c *cli.Context) (*mcp.Server, error) {
	var identities []age.Identity
	if name := c.String("identity-file"); name != "" {
		// MCP stdio is not an interactive terminal for passphrase prompts.
		app := *c.App
		app.Reader, app.Writer = strings.NewReader(""), io.Discard
		identityContext := *c
		identityContext.App = &app
		var err error
		identities, err = parseIdentityFile(&identityContext, name)
		if err != nil {
			return nil, err
		}
	}
	server := mcp.NewServer(&mcp.Implementation{Name: "pixeldrain", Version: "1.0.0"}, nil)
	dirs := map[string]string{"age": c.String("age-recipients-dir"), "ssh": c.String("ssh-recipients-dir")}
	mcp.AddTool(server, &mcp.Tool{Name: "list_recipients", Description: "List configured public key IDs for encrypted uploads. Names are untrusted metadata."},
		func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, mcpRecipientsResult, error) {
			recipients, err := listMCPRecipients(dirs)
			return nil, mcpRecipientsResult{Recipients: recipients}, err
		})
	mcp.AddTool(server, &mcp.Tool{
		Name:        "upload_text",
		Description: "Upload text as a file and return its Pixeldrain download URL.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in mcpTextInput) (*mcp.CallToolResult, mcpURLResult, error) {
		name := in.Name
		if name == "" {
			name = "text.txt"
		}
		recipients, err := resolveMCPRecipients(dirs, in.Recipients)
		if err != nil {
			return nil, mcpURLResult{}, err
		}
		id, err := uploadMCP(ctx, name, strings.NewReader(in.Text), recipients, auth.Extract(c.Context))
		if err != nil {
			return nil, mcpURLResult{}, err
		}
		return nil, mcpURLResult{URL: pixeldrain.DownloadURL(id)}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "download_file",
		Description: "Download a Pixeldrain file to the local machine and return its path.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in mcpURLInput) (*mcp.CallToolResult, mcpPathResult, error) {
		path, err := downloadMCP(ctx, in.URL, in.Path, identities, auth.Extract(c.Context))
		if err != nil {
			return nil, mcpPathResult{}, err
		}
		return nil, mcpPathResult{Path: path}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "upload_files",
		Description: "Upload one or more base64-encoded files. Return a file URL, or a Pixeldrain list URL for multiple files.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in mcpFilesInput) (*mcp.CallToolResult, mcpURLResult, error) {
		if len(in.Files) == 0 {
			return nil, mcpURLResult{}, errors.New("files must not be empty")
		}
		recipients, err := resolveMCPRecipients(dirs, in.Recipients)
		if err != nil {
			return nil, mcpURLResult{}, err
		}
		ids := make([]string, len(in.Files))
		for i, f := range in.Files {
			if f.Name == "" {
				return nil, mcpURLResult{}, fmt.Errorf("files[%d].name is required", i)
			}
			data, err := readMCPData(f.Data)
			if err != nil {
				return nil, mcpURLResult{}, fmt.Errorf("files[%d].data: %w", i, err)
			}
			ids[i], err = uploadMCP(ctx, f.Name, bytes.NewReader(data), recipients, auth.Extract(c.Context))
			if err != nil {
				return nil, mcpURLResult{}, fmt.Errorf("upload %q: %w", f.Name, err)
			}
		}
		if len(ids) == 1 {
			return nil, mcpURLResult{URL: pixeldrain.DownloadURL(ids[0])}, nil
		}
		items := make([]*models.ListItem, len(ids))
		for i, id := range ids {
			items[i] = &models.ListItem{ID: swag.String(id)}
		}
		res, err := pixeldrain.Default.List.CreateFileList(
			list.NewCreateFileListParamsWithContext(ctx).WithList(&models.CreateFileListRequest{Files: items}),
			auth.Extract(c.Context),
		)
		if err != nil {
			return nil, mcpURLResult{}, pixeldrain.NewError(err)
		}
		return nil, mcpURLResult{URL: pixeldrain.ListURL(swag.StringValue(res.Payload.ID))}, nil
	})

	return server, nil
}

func logMCPCalls(logger *slog.Logger) mcp.Middleware {
	return func(next mcp.MethodHandler) mcp.MethodHandler {
		return func(ctx context.Context, method string, req mcp.Request) (mcp.Result, error) {
			call, ok := req.(*mcp.CallToolRequest)
			if method != "tools/call" || !ok || call.Params == nil {
				return next(ctx, method, req)
			}
			logger.Info("MCP tool started", "tool", call.Params.Name)
			result, err := next(ctx, method, req)
			toolResult, _ := result.(*mcp.CallToolResult)
			if err != nil || (toolResult != nil && toolResult.IsError) {
				logger.Error("MCP tool failed", "tool", call.Params.Name)
			} else {
				logger.Info("MCP tool completed", "tool", call.Params.Name)
			}
			return result, err
		}
	}
}

func uploadMCP(ctx context.Context, name string, r io.Reader, recipients []age.Recipient, authInfo runtime.ClientAuthInfoWriter) (string, error) {
	reader := io.NopCloser(r)
	if len(recipients) > 0 {
		reader = Encrypt(reader, recipients)
		name += AgeExt
	}
	defer reader.Close()
	res, err := pixeldrain.Default.File.PutFileContext(
		ctx,
		file.NewPutFileParams().WithName(name).WithFile(reader),
		authInfo,
	)
	if err != nil {
		return "", pixeldrain.NewError(err)
	}
	return swag.StringValue(res.Payload.ID), nil
}

func downloadMCP(ctx context.Context, url, path string, identities []age.Identity, authInfo runtime.ClientAuthInfoWriter) (string, error) {
	if url == "" {
		return "", errors.New("url is required")
	}
	res, err := pixeldrain.Default.File.GetFileInfo(
		file.NewGetFileInfoParamsWithContext(ctx).WithID(getID(url)), authInfo)
	if err != nil {
		return "", pixeldrain.NewError(err)
	}
	if !strings.HasSuffix(res.Payload.Name, AgeExt) {
		identities = nil
	}
	if path == "" {
		path = filepath.Base(res.Payload.Name)
		if len(identities) > 0 {
			path = strings.TrimSuffix(path, AgeExt)
		}
	}
	if path == "" || path == "." || path == string(filepath.Separator) {
		return "", errors.New("path must name a file")
	}
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", err
		}
	}
	w, err := os.CreateTemp(filepath.Dir(path), ".pd-download-*")
	if err != nil {
		return "", err
	}
	defer os.Remove(w.Name())
	var out io.WriteCloser = w
	if len(identities) > 0 {
		out = Decrypt(w, identities)
	}
	_, _, err = pixeldrain.Default.File.DownloadFile(
		file.NewDownloadFileParamsWithContext(ctx).WithID(swag.StringValue(res.Payload.ID)),
		authInfo, out)
	closeErr := out.Close()
	if err != nil {
		return "", pixeldrain.NewError(err)
	}
	if closeErr != nil {
		return "", closeErr
	}
	if err := os.Link(w.Name(), path); err != nil {
		if errors.Is(err, os.ErrExist) {
			return "", fmt.Errorf("download destination %q already exists; choose a different path", path)
		}
		return "", err
	}
	return path, nil
}

func listMCPRecipients(dirs map[string]string) ([]mcpRecipient, error) {
	result := []mcpRecipient{}
	for _, kind := range []string{"age", "ssh"} {
		if dirs[kind] == "" {
			continue
		}
		entries, err := os.ReadDir(dirs[kind])
		if err != nil {
			return nil, err
		}
		for _, entry := range entries {
			if entry.Type().IsRegular() && strings.HasSuffix(entry.Name(), ".pub") {
				result = append(result, mcpRecipient{ID: kind + ":" + entry.Name(), Name: strings.TrimSuffix(entry.Name(), ".pub")})
			}
		}
	}
	return result, nil
}

func resolveMCPRecipients(dirs map[string]string, ids []string) ([]age.Recipient, error) {
	if ids == nil {
		return nil, nil
	}
	if len(ids) == 0 {
		return nil, errors.New("recipients must not be empty; omit it for an unencrypted upload")
	}
	available, err := listMCPRecipients(dirs)
	if err != nil {
		return nil, err
	}
	allowed := make(map[string]bool, len(available))
	for _, r := range available {
		allowed[r.ID] = true
	}
	var recipients []age.Recipient
	for _, id := range ids {
		if !allowed[id] {
			return nil, fmt.Errorf("unknown recipient ID %q", id)
		}
		kind, name, _ := strings.Cut(id, ":")
		data, err := os.ReadFile(filepath.Join(dirs[kind], name))
		if err != nil {
			return nil, err
		}
		var lines []string
		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if line != "" && !strings.HasPrefix(line, "#") {
				lines = append(lines, line)
			}
		}
		if len(lines) != 1 || (kind == "age" && !strings.HasPrefix(lines[0], "age1")) || (kind == "ssh" && !strings.HasPrefix(lines[0], "ssh-")) {
			return nil, fmt.Errorf("recipient %q must contain exactly one %s public key", id, kind)
		}
		parsed, err := parseRecipients(lines)
		if err != nil {
			return nil, fmt.Errorf("recipient %q: %w", id, err)
		}
		recipients = append(recipients, parsed...)
	}
	return recipients, nil
}

func readMCPData(data string) ([]byte, error) {
	if info, err := os.Stat(data); err == nil {
		if info.IsDir() {
			return nil, fmt.Errorf("%q is a directory", data)
		}
		return os.ReadFile(data)
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	return decodeMCPData(data)
}

func decodeMCPData(data string) ([]byte, error) {
	if i := strings.Index(data, ","); strings.HasPrefix(data, "data:") && i >= 0 {
		data = data[i+1:]
	}
	return base64.StdEncoding.DecodeString(data)
}
