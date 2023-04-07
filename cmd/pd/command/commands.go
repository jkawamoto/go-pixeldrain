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
	// FlagRecipient is the flag to specify recipient's public key.
	FlagRecipient = "recipient"
	// FlagRecipientFile is the flag to specify a path to the file that contains recipients' public keys.
	FlagRecipientFile = "recipient-file"
	// FlagIdentity is the flag to specify a path to the file that contains the user's private keys.
	FlagIdentity = "identity"
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
					"This flag can specifies the album `NAME`.",
			},
			&cli.StringSliceFlag{
				Name:     FlagRecipient,
				Aliases:  []string{"r"},
				Category: "End-to-end encryption",
				Usage:    "Encrypt to the specified `RECIPIENT`. Can be repeated.",
			},
			&cli.StringFlag{
				Name:     FlagRecipientFile,
				Aliases:  []string{"R"},
				Category: "End-to-end encryption",
				Usage:    "Encrypt to recipients listed at `PATH`.",
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
			&cli.StringFlag{
				Name:     FlagIdentity,
				Aliases:  []string{"i"},
				Category: "End-to-end encryption",
				Usage:    "Use the identity file at `PATH`.",
			},
		},
	},
}
