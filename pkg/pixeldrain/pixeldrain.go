// pixeldrain.go
//
// Copyright (c) 2018-2021 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package pixeldrain

import (
	"io"
	"os"
	"path"

	"github.com/go-openapi/runtime"
	auth "github.com/go-openapi/runtime/client"

	"github.com/jkawamoto/go-pixeldrain/pkg/pixeldrain/client"
)

type Pixeldrain struct {
	cli            *client.PixeldrainAPI
	authInfoWriter runtime.ClientAuthInfoWriter
	Stdout         io.Writer
	Stderr         io.Writer
}

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

func (Pixeldrain) DownloadURL(id string) string {
	return "https://" + path.Join(client.DefaultHost, client.DefaultBasePath, "file", id)
}
