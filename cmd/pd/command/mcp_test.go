package command

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
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
