package command

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"flag"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"filippo.io/age"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"
	"github.com/golang/mock/gomock"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/ssh"

	"github.com/jkawamoto/go-pixeldrain/client/file"
	"github.com/jkawamoto/go-pixeldrain/cmd/pd/command/mock"
	"github.com/jkawamoto/go-pixeldrain/models"
)

func TestLogMCPCalls(t *testing.T) {
	for _, test := range []struct {
		name   string
		result mcp.Result
		err    error
		want   string
	}{
		{name: "success", result: &mcp.CallToolResult{}, want: "MCP tool completed"},
		{name: "tool error", result: &mcp.CallToolResult{IsError: true}, want: "MCP tool failed"},
		{name: "protocol error", err: errors.New("secret error"), want: "MCP tool failed"},
	} {
		t.Run(test.name, func(t *testing.T) {
			var logs bytes.Buffer
			logger := slog.New(slog.NewTextHandler(&logs, nil))
			handler := logMCPCalls(logger)(func(context.Context, string, mcp.Request) (mcp.Result, error) {
				return test.result, test.err
			})
			req := &mcp.CallToolRequest{Params: &mcp.CallToolParamsRaw{Name: "upload_text", Arguments: []byte(`{"text":"secret"}`)}}
			result, err := handler(context.Background(), "tools/call", req)
			if result != test.result || err != test.err {
				t.Fatal("middleware changed the result or error")
			}
			if got := logs.String(); !strings.Contains(got, "MCP tool started") || !strings.Contains(got, test.want) || !strings.Contains(got, "upload_text") || strings.Contains(got, "secret") {
				t.Fatalf("unexpected logs: %s", got)
			}
		})
	}
}

func TestMCPToolSchemas(t *testing.T) {
	server := mcp.NewServer(&mcp.Implementation{Name: "test"}, nil)
	mcp.AddTool(server, &mcp.Tool{Name: "upload_text"}, func(context.Context, *mcp.CallToolRequest, mcpTextInput) (*mcp.CallToolResult, mcpURLResult, error) {
		return nil, mcpURLResult{}, nil
	})
}

func TestMCPRecipientTools(t *testing.T) {
	dir := t.TempDir()
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "alice.pub"), []byte(identity.Recipient().String()), 0600); err != nil {
		t.Fatal(err)
	}
	flags := flag.NewFlagSet("mcp", flag.ContinueOnError)
	flags.String("age-recipients-dir", dir, "")
	c := cli.NewContext(cli.NewApp(), flags, nil)
	ctx := context.Background()
	ct, st := mcp.NewInMemoryTransports()
	server, err := newMCPServer(c)
	if err != nil {
		t.Fatal(err)
	}
	serverSession, err := server.Connect(ctx, st, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer serverSession.Close()
	client, err := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "1"}, nil).Connect(ctx, ct, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	result, err := client.CallTool(ctx, &mcp.CallToolParams{Name: "list_recipients", Arguments: map[string]any{}})
	if err != nil || result.IsError {
		t.Fatalf("list: %v %v", result, err)
	}
	m := mock.NewMockClientService(gomock.NewController(t))
	RegisterMock(t, m)
	m.EXPECT().PutFileContext(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, params *file.PutFileParams, _ runtime.ClientAuthInfoWriter, _ ...file.ClientOption) (*file.PutFileCreated, error) {
			r, err := age.Decrypt(params.File, identity)
			if err != nil {
				return nil, err
			}
			data, err := io.ReadAll(r)
			if string(data) != "secret" {
				t.Errorf("plaintext = %q", data)
			}
			return &file.PutFileCreated{Payload: &models.SuccessResponse{ID: swag.String("id")}}, err
		})
	result, err = client.CallTool(ctx, &mcp.CallToolParams{Name: "upload_text", Arguments: map[string]any{"text": "secret", "recipients": []string{"age:alice.pub"}}})
	if err != nil || result.IsError {
		t.Fatalf("upload: %v %v", result, err)
	}
	result, err = client.CallTool(ctx, &mcp.CallToolParams{Name: "upload_text", Arguments: map[string]any{"text": "secret", "recipients": []string{}}})
	if err != nil || !result.IsError {
		t.Fatalf("empty recipients: %v %v", result, err)
	}
}

func TestMCPConfiguredIdentity(t *testing.T) {
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	keyPath := filepath.Join(dir, "identity")
	if err := os.WriteFile(keyPath, []byte(identity.String()), 0600); err != nil {
		t.Fatal(err)
	}
	flags := flag.NewFlagSet("mcp", flag.ContinueOnError)
	flags.String("identity-file", keyPath, "")
	c := cli.NewContext(cli.NewApp(), flags, nil)
	server, err := newMCPServer(c)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	ct, st := mcp.NewInMemoryTransports()
	ss, err := server.Connect(ctx, st, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer ss.Close()
	client, err := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "1"}, nil).Connect(ctx, ct, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	var encrypted bytes.Buffer
	w, err := age.Encrypt(&encrypted, identity.Recipient())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := w.Write([]byte("secret")); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	m := mock.NewMockClientService(gomock.NewController(t))
	RegisterMock(t, m)
	for _, tc := range []struct{ name, content string }{{"secret.age", encrypted.String()}, {"plain.txt", "plain"}} {
		m.EXPECT().GetFileInfo(gomock.Any(), gomock.Any()).Return(&file.GetFileInfoOK{Payload: &models.FileInfo{ID: swag.String("id"), Name: tc.name}}, nil)
		m.EXPECT().DownloadFile(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ *file.DownloadFileParams, _ runtime.ClientAuthInfoWriter, out io.Writer, _ ...file.ClientOption) (*file.DownloadFileOK, *file.DownloadFilePartialContent, error) {
				_, err := io.WriteString(out, tc.content)
				return &file.DownloadFileOK{}, nil, err
			})
		path := filepath.Join(dir, tc.name+".out")
		result, err := client.CallTool(ctx, &mcp.CallToolParams{Name: "download_file", Arguments: map[string]any{"url": "https://pixeldrain.com/u/id", "path": path}})
		if err != nil || result.IsError {
			t.Fatalf("download %s: %v %v", tc.name, result, err)
		}
		want := "secret"
		if tc.name == "plain.txt" {
			want = "plain"
		}
		if data, err := os.ReadFile(path); err != nil || string(data) != want {
			t.Fatalf("download %s: %q %v", tc.name, data, err)
		}
	}
	if err := flags.Set("identity-file", filepath.Join(dir, "missing")); err != nil {
		t.Fatal(err)
	}
	if _, err := newMCPServer(c); err == nil {
		t.Fatal("missing identity accepted")
	}
}

func TestMCPRecipients(t *testing.T) {
	ageDir, sshDir := t.TempDir(), t.TempDir()
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatal(err)
	}
	pub, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	sshKey, err := ssh.NewPublicKey(pub)
	if err != nil {
		t.Fatal(err)
	}
	for path, data := range map[string]string{
		filepath.Join(ageDir, "alice.pub"): "# age key\n" + identity.Recipient().String() + "\n",
		filepath.Join(sshDir, "alice.pub"): string(ssh.MarshalAuthorizedKey(sshKey)),
		filepath.Join(sshDir, "private"):   "ignored",
	} {
		if err := os.WriteFile(path, []byte(data), 0600); err != nil {
			t.Fatal(err)
		}
	}
	dirs := map[string]string{"age": ageDir, "ssh": sshDir}
	got, err := listMCPRecipients(dirs)
	if err != nil {
		t.Fatal(err)
	}
	want := []mcpRecipient{{ID: "age:alice.pub", Name: "alice"}, {ID: "ssh:alice.pub", Name: "alice"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	recipients, err := resolveMCPRecipients(dirs, []string{"age:alice.pub", "ssh:alice.pub"})
	if err != nil || len(recipients) != 2 {
		t.Fatalf("recipients: %v, %v", recipients, err)
	}
	for _, ids := range [][]string{{}, {"age:../alice.pub"}, {"ssh:missing.pub"}} {
		if _, err := resolveMCPRecipients(dirs, ids); err == nil {
			t.Errorf("accepted %v", ids)
		}
	}
	if err := os.WriteFile(filepath.Join(ageDir, "alice.pub"), []byte("not a key"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := resolveMCPRecipients(dirs, []string{"age:alice.pub"}); err == nil {
		t.Fatal("accepted invalid key")
	}
}

func TestMCPEncryptedTransfer(t *testing.T) {
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatal(err)
	}
	wrong, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatal(err)
	}
	m := mock.NewMockClientService(gomock.NewController(t))
	RegisterMock(t, m)
	var ciphertext []byte
	m.EXPECT().PutFileContext(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, params *file.PutFileParams, _ runtime.ClientAuthInfoWriter, _ ...file.ClientOption) (*file.PutFileCreated, error) {
			if params.Name != "secret.txt.age" {
				t.Errorf("name = %q", params.Name)
			}
			var err error
			ciphertext, err = io.ReadAll(params.File)
			return &file.PutFileCreated{Payload: &models.SuccessResponse{ID: swag.String("id")}}, err
		})
	if _, err := uploadMCP(context.Background(), "secret.txt", strings.NewReader("secret"), []age.Recipient{identity.Recipient()}, nil); err != nil {
		t.Fatal(err)
	}
	m.EXPECT().GetFileInfo(gomock.Any(), gomock.Any()).Return(&file.GetFileInfoOK{Payload: &models.FileInfo{ID: swag.String("id"), Name: "secret.txt.age"}}, nil).Times(3)
	m.EXPECT().DownloadFile(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ *file.DownloadFileParams, _ runtime.ClientAuthInfoWriter, w io.Writer, _ ...file.ClientOption) (*file.DownloadFileOK, *file.DownloadFilePartialContent, error) {
			_, err := w.Write(ciphertext)
			return &file.DownloadFileOK{}, nil, err
		}).Times(3)
	path := filepath.Join(t.TempDir(), "output")
	if err := os.WriteFile(path, []byte("existing"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := downloadMCP(context.Background(), "https://pixeldrain.com/u/id", path, []age.Identity{wrong}, nil); err == nil {
		t.Fatal("wrong key accepted")
	}
	if data, err := os.ReadFile(path); err != nil || string(data) != "existing" {
		t.Fatalf("existing file changed: %q, %v", data, err)
	}
	if _, err := downloadMCP(context.Background(), "https://pixeldrain.com/u/id", path, []age.Identity{identity}, nil); err == nil {
		t.Fatal("existing file overwritten")
	}
	if data, err := os.ReadFile(path); err != nil || string(data) != "existing" {
		t.Fatalf("existing file changed: %q, %v", data, err)
	}
	path = filepath.Join(filepath.Dir(path), "new-output")
	if _, err := downloadMCP(context.Background(), "https://pixeldrain.com/u/id", path, []age.Identity{identity}, nil); err != nil {
		t.Fatal(err)
	}
	if data, err := os.ReadFile(path); err != nil || string(data) != "secret" {
		t.Fatalf("download: %q, %v", data, err)
	}
}

func TestDecodeMCPData(t *testing.T) {
	for _, test := range []struct {
		name string
		data string
		want []byte
	}{
		{name: "base64", data: "aGVsbG8=", want: []byte("hello")},
		{name: "data URI", data: "data:text/plain;base64,aGVsbG8=", want: []byte("hello")},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := decodeMCPData(test.data)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(got, test.want) {
				t.Fatalf("decodeMCPData() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestReadMCPDataFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "input.txt")
	if err := os.WriteFile(path, []byte("hello"), 0600); err != nil {
		t.Fatal(err)
	}

	got, err := readMCPData(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, []byte("hello")) {
		t.Fatalf("readMCPData() = %q, want %q", got, "hello")
	}
}
