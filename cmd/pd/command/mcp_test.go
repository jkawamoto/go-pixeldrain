package command

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"filippo.io/age"
	"filippo.io/age/agessh"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"
	"github.com/golang/mock/gomock"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/ssh"

	"github.com/jkawamoto/go-pixeldrain/client/file"
	"github.com/jkawamoto/go-pixeldrain/client/list"
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
	flags := flag.NewFlagSet("mcp", flag.ContinueOnError)
	flags.String("age-recipients-dir", "", "")
	flags.String("ssh-recipients-dir", "", "")
	c := cli.NewContext(cli.NewApp(), flags, nil)
	server, err := newMCPServer(c)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	serverSession, err := server.Connect(ctx, serverTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer serverSession.Close()
	client, err := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "1"}, nil).Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	result, err := client.ListTools(ctx, &mcp.ListToolsParams{})
	if err != nil {
		t.Fatal(err)
	}
	tools := make(map[string]*mcp.Tool, len(result.Tools))
	for _, tool := range result.Tools {
		tools[tool.Name] = tool
	}
	for _, name := range []string{"list_recipients", "upload_text", "upload_files"} {
		if tools[name] == nil {
			t.Fatalf("tools/list did not return %q", name)
		}
	}
	for _, phrase := range []string{"github:USERNAME", "not listed", "untrusted metadata"} {
		if !strings.Contains(tools["list_recipients"].Description, phrase) {
			t.Errorf("list_recipients description %q does not contain %q", tools["list_recipients"].Description, phrase)
		}
	}
	for _, name := range []string{"upload_text", "upload_files"} {
		t.Run(name, func(t *testing.T) {
			tool := tools[name]
			for _, phrase := range []string{"github:octocat", "do not appear", "matching private key", "trusted channel", "never falls back to plaintext", "warnings"} {
				if !strings.Contains(tool.Description, phrase) {
					t.Errorf("description %q does not contain %q", tool.Description, phrase)
				}
			}
			recipientsDescription := schemaPropertyDescription(t, tool.InputSchema, "recipients")
			for _, phrase := range []string{"github:octocat", "empty array is an error", "plaintext fallback", "trusted channel", "warnings"} {
				if !strings.Contains(recipientsDescription, phrase) {
					t.Errorf("recipients schema description %q does not contain %q", recipientsDescription, phrase)
				}
			}
			for _, property := range []string{"url", "warnings"} {
				if description := schemaPropertyDescription(t, tool.OutputSchema, property); description == "" {
					t.Errorf("output schema property %q has no description", property)
				}
			}
		})
	}
	if description := schemaPropertyDescription(t, tools["list_recipients"].OutputSchema, "recipients"); description == "" {
		t.Error("list_recipients output schema does not describe recipients")
	}
}

func schemaPropertyDescription(t *testing.T, schema any, property string) string {
	t.Helper()
	schemaMap, ok := schema.(map[string]any)
	if !ok {
		t.Fatalf("schema has type %T, want map[string]any", schema)
	}
	properties, ok := schemaMap["properties"].(map[string]any)
	if !ok {
		t.Fatalf("schema properties have type %T, want map[string]any", schemaMap["properties"])
	}
	propertySchema, ok := properties[property].(map[string]any)
	if !ok {
		t.Fatalf("schema property %q has type %T, want map[string]any", property, properties[property])
	}
	description, _ := propertySchema["description"].(string)
	return description
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

func TestMCPGitHubRecipients(t *testing.T) {
	localIdentity, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatal(err)
	}
	sshPublic, sshPrivate, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	sshPublicKey, err := ssh.NewPublicKey(sshPublic)
	if err != nil {
		t.Fatal(err)
	}
	sshIdentity, err := agessh.NewEd25519Identity(sshPrivate)
	if err != nil {
		t.Fatal(err)
	}
	ecdsaKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	unsupportedKey, err := ssh.NewPublicKey(&ecdsaKey.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	githubKeys := string(ssh.MarshalAuthorizedKey(sshPublicKey)) + string(ssh.MarshalAuthorizedKey(unsupportedKey))
	githubRequests := 0
	httpClient := &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		githubRequests++
		if request.URL.String() != "https://github.com/octocat.keys" {
			t.Errorf("unexpected GitHub request URL: %s", request.URL)
		}
		return githubResponse(http.StatusOK, githubKeys), nil
	})}

	ageDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(ageDir, "alice.pub"), []byte(localIdentity.Recipient().String()), 0600); err != nil {
		t.Fatal(err)
	}
	flags := flag.NewFlagSet("mcp", flag.ContinueOnError)
	flags.String("age-recipients-dir", ageDir, "")
	flags.String("ssh-recipients-dir", "", "")
	c := cli.NewContext(cli.NewApp(), flags, nil)
	ctx := context.Background()
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	server, err := newMCPServerWithHTTPClient(c, httpClient)
	if err != nil {
		t.Fatal(err)
	}
	serverSession, err := server.Connect(ctx, serverTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer serverSession.Close()
	client, err := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "1"}, nil).Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	m := mock.NewMockClientService(gomock.NewController(t))
	RegisterMock(t, m)
	wantPlaintext := map[string]string{"text.txt.age": "secret", "first.txt.age": "first", "second.txt.age": "second"}
	m.EXPECT().PutFileContext(gomock.Any(), gomock.Any(), gomock.Any()).Times(3).DoAndReturn(
		func(_ context.Context, params *file.PutFileParams, _ runtime.ClientAuthInfoWriter, _ ...file.ClientOption) (*file.PutFileCreated, error) {
			if !strings.HasSuffix(params.Name, AgeExt) {
				t.Errorf("uploaded name = %q, want .age suffix", params.Name)
			}
			ciphertext, err := io.ReadAll(params.File)
			if err != nil {
				return nil, err
			}
			for _, identity := range []age.Identity{localIdentity, sshIdentity} {
				reader, err := age.Decrypt(bytes.NewReader(ciphertext), identity)
				if err != nil {
					t.Errorf("decrypt with recipient identity: %v", err)
					continue
				}
				plaintext, err := io.ReadAll(reader)
				if err != nil || string(plaintext) != wantPlaintext[params.Name] {
					t.Errorf("decrypted %q = %q, %v", params.Name, plaintext, err)
				}
			}
			return &file.PutFileCreated{Payload: &models.SuccessResponse{ID: swag.String("uploaded")}}, nil
		})
	m.EXPECT().CreateFileList(gomock.Any(), gomock.Any(), gomock.Any()).Return(
		&list.CreateFileListCreated{Payload: &models.SuccessResponse{ID: swag.String("album")}}, nil)

	recipients := []string{"github:octocat", "age:alice.pub"}
	result, err := client.CallTool(ctx, &mcp.CallToolParams{
		Name:      "upload_text",
		Arguments: map[string]any{"text": "secret", "recipients": recipients},
	})
	if err != nil || result.IsError {
		t.Fatalf("upload_text: %v %v", result, err)
	}
	gotText := decodeMCPURLResult(t, result)
	if len(gotText.Warnings) != 1 || !strings.Contains(gotText.Warnings[0], "unsupported SSH key type") {
		t.Fatalf("upload_text warnings = %v", gotText.Warnings)
	}

	result, err = client.CallTool(ctx, &mcp.CallToolParams{
		Name: "upload_files",
		Arguments: map[string]any{
			"recipients": recipients,
			"files": []map[string]string{
				{"name": "first.txt", "data": base64.StdEncoding.EncodeToString([]byte("first"))},
				{"name": "second.txt", "data": base64.StdEncoding.EncodeToString([]byte("second"))},
			},
		},
	})
	if err != nil || result.IsError {
		t.Fatalf("upload_files: %v %v", result, err)
	}
	gotFiles := decodeMCPURLResult(t, result)
	if len(gotFiles.Warnings) != 1 || githubRequests != 2 {
		t.Fatalf("upload_files warnings = %v; GitHub requests = %d", gotFiles.Warnings, githubRequests)
	}
}

func TestMCPGitHubRecipientFailuresDoNotUpload(t *testing.T) {
	for _, tc := range []struct {
		name      string
		status    int
		body      string
		identity  string
		wantError string
	}{
		{name: "GitHub fetch failure", status: http.StatusNotFound, wantError: "HTTP status 404"},
		{name: "hybrid mix", status: http.StatusOK, body: "", identity: "hybrid", wantError: "hybrid age keys"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ageDir := t.TempDir()
			if tc.identity == "hybrid" {
				identity, err := age.GenerateHybridIdentity()
				if err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(ageDir, "hybrid.pub"), []byte(identity.Recipient().String()), 0600); err != nil {
					t.Fatal(err)
				}
			}
			flags := flag.NewFlagSet("mcp", flag.ContinueOnError)
			flags.String("age-recipients-dir", ageDir, "")
			flags.String("ssh-recipients-dir", "", "")
			c := cli.NewContext(cli.NewApp(), flags, nil)
			httpClient := githubClient(tc.status, tc.body, nil)
			if tc.identity == "hybrid" {
				publicKey, _, err := ed25519.GenerateKey(rand.Reader)
				if err != nil {
					t.Fatal(err)
				}
				sshPublicKey, err := ssh.NewPublicKey(publicKey)
				if err != nil {
					t.Fatal(err)
				}
				httpClient = githubClient(tc.status, string(ssh.MarshalAuthorizedKey(sshPublicKey)), nil)
			}
			server, err := newMCPServerWithHTTPClient(c, httpClient)
			if err != nil {
				t.Fatal(err)
			}
			ctx := context.Background()
			clientTransport, serverTransport := mcp.NewInMemoryTransports()
			serverSession, err := server.Connect(ctx, serverTransport, nil)
			if err != nil {
				t.Fatal(err)
			}
			defer serverSession.Close()
			client, err := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "1"}, nil).Connect(ctx, clientTransport, nil)
			if err != nil {
				t.Fatal(err)
			}
			defer client.Close()

			RegisterMock(t, mock.NewMockClientService(gomock.NewController(t)))
			recipients := []string{"github:octocat"}
			if tc.identity == "hybrid" {
				recipients = append(recipients, "age:hybrid.pub")
			}
			result, err := client.CallTool(ctx, &mcp.CallToolParams{
				Name:      "upload_text",
				Arguments: map[string]any{"text": "secret", "recipients": recipients},
			})
			if err != nil || !result.IsError || !mcpResultContains(result, tc.wantError) {
				t.Fatalf("expected error containing %q, got result=%v err=%v", tc.wantError, result, err)
			}
		})
	}
}

func decodeMCPURLResult(t *testing.T, result *mcp.CallToolResult) mcpURLResult {
	t.Helper()
	data, err := json.Marshal(result.StructuredContent)
	if err != nil {
		t.Fatal(err)
	}
	var decoded mcpURLResult
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	return decoded
}

func mcpResultContains(result *mcp.CallToolResult, text string) bool {
	for _, content := range result.Content {
		if textContent, ok := content.(*mcp.TextContent); ok && strings.Contains(textContent.Text, text) {
			return true
		}
	}
	return false
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
