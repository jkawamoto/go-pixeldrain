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

	"github.com/jkawamoto/go-pixeldrain/pkg/pixeldrain/client"
)

var DownloadEndpoint = "https://" + path.Join(client.DefaultHost, client.DefaultBasePath, "file")

type Pixeldrain struct {
	Client *client.PixelDrain
	Stdout io.Writer
	Stderr io.Writer
}

func New() *Pixeldrain {
	cli := client.Default
	cli.SetTransport(newTransport(cli.Transport))

	return &Pixeldrain{
		Client: cli,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}
