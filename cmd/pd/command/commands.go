// commands.go
//
// Copyright (c) 2018-2023 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package command

import (
	"errors"

	"github.com/urfave/cli/v2"
)

const (
	FlagAlbumName = "album"
)

var (
	ErrNotEnoughArguments = errors.New("not enough arguments")
)

// Commands manage sub commands.
var Commands = []*cli.Command{
	{
		Name:  "upload",
		Usage: "Upload files",
		Description: "Upload files specified by the given paths. Each path can have an optional name. " +
			"If a name is given, uploaded file will be renamed with it.",
		ArgsUsage: "<path[:name]>...",
		Action:    CmdUpload,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: FlagAlbumName,
				Usage: "If multiple files are uploaded, an album consisting of those files will be created. " +
					"This flag can specifies the album name.",
			},
		},
	}, {
		Name:        "download",
		Usage:       "Download a file",
		Description: "download a file from PixelDrain",
		ArgsUsage:   "<file ID | URL>",
		Action:      CmdDownload,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "output the downloaded file into `DIR`",
			},
		},
	},
}
