/*
 * commands.go
 *
 * Copyright (c) 2018-2019 Junpei Kawamoto
 *
 * This software is released under the MIT License.
 *
 * http://opensource.org/licenses/mit-license.php
 */

package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"

	"github.com/jkawamoto/go-pixeldrain/cmd/pd/command"
	"github.com/jkawamoto/go-pixeldrain/cmd/pd/status"
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
	}, {
		Name:        "download",
		Usage:       "Download a file",
		Description: "download a file from PixelDrain",
		ArgsUsage:   "<file ID | URL>",
		Action:      command.CmdDownload,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "o, output",
				Usage: "output the downloaded file into `DIR`",
			},
		},
	}, {
		Name:        "create-list",
		Usage:       "Create a list consisting of uploaded files",
		Description: "create a list consisting of given file IDs",
		ArgsUsage:   "fileID[:description]...",
		Action:      command.CmdCreateList,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "t, title",
				Usage: "specify the `TITLE` of this list",
			},
			cli.StringFlag{
				Name:  "description",
				Usage: "specify the description of this list",
			},
		},
	},
}

// CommandNotFound shows error message and exit when a given command is not found.
func CommandNotFound(c *cli.Context, command string) {
	_, _ = fmt.Fprintf(os.Stderr, "'%s' is not a %s command..\n", command, c.App.Name)
	_ = cli.ShowAppHelp(c)
	os.Exit(status.InvalidCommand)
}
