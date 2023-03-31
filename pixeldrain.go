// pixeldrain.go
//
// Copyright (c) 2018-2023 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package pixeldrain

import (
	"io"
	"net/url"
	"os"

	"github.com/go-openapi/runtime"
	auth "github.com/go-openapi/runtime/client"

	"github.com/jkawamoto/go-pixeldrain/client"
)

// Pixeldrain is a Pixeldrain API client.
type Pixeldrain struct {
	// Stdout is used to output downloaded files.
	Stdout io.Writer
	// Stderr is used to render progress bars. If you want to disable progress bars, set io.Discard.
	Stderr io.Writer

	cli            *client.PixeldrainAPI
	authInfoWriter runtime.ClientAuthInfoWriter
}

// New creates a Pixeldrain API client that uses the given API key. The key can be an empty string.
func New(apiKey string) *Pixeldrain {
	cli := client.Default
	cli.SetTransport(newTransport(cli.Transport))

	res := &Pixeldrain{
		cli:    cli,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	if apiKey != "" {
		res.authInfoWriter = auth.BasicAuth("", apiKey)
	}
	return res
}

// DownloadURL returns the URL associated with the given file ID.
func (*Pixeldrain) DownloadURL(id string) string {
	u := url.URL{
		Scheme: client.DefaultSchemes[0],
	}
	return u.JoinPath(client.DefaultHost, client.DefaultBasePath, "file", id).String()
}

// ListURL returns the URL associated with the given list ID.
func (*Pixeldrain) ListURL(id string) string {
	u := url.URL{
		Scheme: client.DefaultSchemes[0],
	}
	return u.JoinPath(client.DefaultHost, "l", id).String()
}
