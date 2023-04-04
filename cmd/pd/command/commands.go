// commands.go
//
// Copyright (c) 2018-2023 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package command

import (
	"github.com/urfave/cli/v2"
)

const (
	// FlagAlbumName is the flag name to specify an album name.
	FlagAlbumName = "album"
	// FlagDirectory is the flag name to specify an output directory where downloaded files will be stored.
	FlagDirectory = "dir"
	// FlagAll is the flag to download all files in a list.
	FlagAll = "all"
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
					"This flag can specifies the album `name`.",
			},
		},
	}, {
		Name:  "download",
		Usage: "Download files",
		Description: "Download files associated with the given URLs. " +
			"If the URL refers an album, you will be asked which file you want to download.",
		ArgsUsage: "<URL>...",
		Action:    CmdDownload,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        FlagDirectory,
				Aliases:     []string{"o"},
				Usage:       "`path` to the directory where downloaded files will be stored",
				DefaultText: ".",
			},
			&cli.BoolFlag{
				Name:  FlagAll,
				Usage: "if an album URL is given, download all files in it",
			},
		},
	},
}
