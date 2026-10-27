package command

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

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
