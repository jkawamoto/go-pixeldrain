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
	"github.com/jkawamoto/go-pixeldrain/client"
	"io"
	"os"
)

type Pixeldrain struct {
	Client *client.PixelDrain
	Stderr io.Writer
}

func New() *Pixeldrain {
	return &Pixeldrain{
		Client: client.Default,
		Stderr: os.Stderr,
	}
}
