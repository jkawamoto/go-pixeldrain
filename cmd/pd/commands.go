// Copyright (c) 2018 Junpei Kawamoto
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"fmt"
	"os"

	"github.com/jkawamoto/go-pixeldrain/command"
	"github.com/urfave/cli"
)

// GlobalFlags manages global flags.
var GlobalFlags []cli.Flag

// Commands manage sub commands.
var Commands = []cli.Command{
	{
		Name:        "upload",
		Usage:       "Upload a file",
		Description: "upload a file to PixelDrain",
		ArgsUsage:   "<file path>",
		Action:      command.CmdUpload,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "n, name",
				Usage: "rename the uploaded file to `NAME`",
			},
		},
	},
}

// CommandNotFound shows error message and exit when a given command is not found.
func CommandNotFound(c *cli.Context, command string) {

	fmt.Fprintf(os.Stderr, "'%s' is not a %s command..\n", command, c.App.Name)
	cli.ShowAppHelp(c)
	os.Exit(2)

}
