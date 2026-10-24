// mcp.go
//
// Copyright (c) 2018-2025 Junpei Kawamoto
//
// This software is released under the MIT License.

package command

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/urfave/cli/v2"

	"github.com/jkawamoto/go-pixeldrain"
	"github.com/jkawamoto/go-pixeldrain/client/file"
	"github.com/jkawamoto/go-pixeldrain/client/list"
	"github.com/jkawamoto/go-pixeldrain/cmd/pd/auth"
	"github.com/jkawamoto/go-pixeldrain/models"
)

type mcpTextInput struct {
	Text string `json:"text" jsonschema:"Text to upload"`
	Name string `json:"name,omitempty" jsonschema:"File name (defaults to text.txt)"`
}

type mcpFile struct {
	Name string `json:"name" jsonschema:"File name"`
	Data string `json:"data" jsonschema:"Base64 encoded file data (data URI is also accepted)"`
}

type mcpFilesInput struct {
	Files []mcpFile `json:"files" jsonschema:"Files to upload"`
}

type mcpURLInput struct {
	URL  string `json:"url" jsonschema:"Pixeldrain file download URL"`
	Path string `json:"path,omitempty" jsonschema:"Local output path; defaults to the downloaded file name"`
}

type mcpURLResult struct {
	URL string `json:"url"`
}

type mcpPathResult struct {
	Path string `json:"path"`
}

func CmdMCP(c *cli.Context) error {
	server := mcp.NewServer(&mcp.Implementation{Name: "pixeldrain", Version: "1.0.0"}, nil)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "upload_text",
		Description: "Upload text as a file and return its Pixeldrain download URL.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in mcpTextInput) (*mcp.CallToolResult, mcpURLResult, error) {
		name := in.Name
		if name == "" {
			name = "text.txt"
		}
		id, err := uploadMCP(ctx, name, strings.NewReader(in.Text), auth.Extract(c.Context))
		if err != nil {
			return nil, mcpURLResult{}, err
		}
		return nil, mcpURLResult{URL: pixeldrain.DownloadURL(id)}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "download_file",
		Description: "Download a Pixeldrain file to the local machine and return its path.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in mcpURLInput) (*mcp.CallToolResult, mcpPathResult, error) {
		path, err := downloadMCP(ctx, in.URL, in.Path, auth.Extract(c.Context))
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
		ids := make([]string, len(in.Files))
		for i, f := range in.Files {
			if f.Name == "" {
				return nil, mcpURLResult{}, fmt.Errorf("files[%d].name is required", i)
			}
			data, err := decodeMCPData(f.Data)
			if err != nil {
				return nil, mcpURLResult{}, fmt.Errorf("files[%d].data: %w", i, err)
			}
			ids[i], err = uploadMCP(ctx, f.Name, bytes.NewReader(data), auth.Extract(c.Context))
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

	return server.Run(context.Background(), &mcp.StdioTransport{})
}

func uploadMCP(ctx context.Context, name string, r io.Reader, authInfo runtime.ClientAuthInfoWriter) (string, error) {
	res, err := pixeldrain.Default.File.UploadFile(
		file.NewUploadFileParamsWithContext(ctx).WithFile(runtime.NamedReader(name, io.NopCloser(r))),
		authInfo,
	)
	if err != nil {
		return "", pixeldrain.NewError(err)
	}
	return swag.StringValue(res.Payload.ID), nil
}

func downloadMCP(ctx context.Context, url, path string, authInfo runtime.ClientAuthInfoWriter) (string, error) {
	if url == "" {
		return "", errors.New("url is required")
	}
	res, err := pixeldrain.Default.File.GetFileInfo(
		file.NewGetFileInfoParamsWithContext(ctx).WithID(getID(url)), authInfo)
	if err != nil {
		return "", pixeldrain.NewError(err)
	}
	if path == "" {
		path = filepath.Base(res.Payload.Name)
	}
	if path == "." || path == string(filepath.Separator) {
		return "", errors.New("path must name a file")
	}
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", err
		}
	}
	w, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer func() { _ = w.Close() }()
	_, _, err = pixeldrain.Default.File.DownloadFile(
		file.NewDownloadFileParamsWithContext(ctx).WithID(swag.StringValue(res.Payload.ID)),
		authInfo, w)
	if err != nil {
		return "", pixeldrain.NewError(err)
	}
	return path, nil
}

func decodeMCPData(data string) ([]byte, error) {
	if i := strings.Index(data, ","); strings.HasPrefix(data, "data:") && i >= 0 {
		data = data[i+1:]
	}
	return base64.StdEncoding.DecodeString(data)
}
