/*
 * pixeldrain.go
 *
 * Copyright (c) 2018-2019 Junpei Kawamoto
 *
 * This software is released under the MIT License.
 *
 * http://opensource.org/licenses/mit-license.php
 */

package pixeldrain

import (
	httpTransport "github.com/go-openapi/runtime/client"
	"github.com/jkawamoto/go-pixeldrain/client"
	"io"
	"os"
)

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
