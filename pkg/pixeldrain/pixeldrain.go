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

	httpTransport "github.com/go-openapi/runtime/client"

	"github.com/jkawamoto/go-pixeldrain/pkg/pixeldrain/client"
)

var DownloadEndpoint = "https://" + path.Join(client.DefaultHost, client.DefaultBasePath, "file")

type Pixeldrain struct {
	Client *client.PixelDrain
	Stdout io.WriteCloser
	Stderr io.WriteCloser
}

func New() *Pixeldrain {
	cli := client.Default

	switch transport := cli.Transport.(type) {
	case *httpTransport.Runtime:
		transport.Transport = newTransporter(transport.Transport)
		cli.SetTransport(transport)
	}

	return &Pixeldrain{
		Client: cli,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}
