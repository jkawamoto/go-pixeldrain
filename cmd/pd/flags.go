// flags.go
//
// Copyright (c) 2018-2025 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package main

import "github.com/urfave/cli/v2"

const (
	FlagAPIKey          = "api-key"
	EnvPixeldrainApiKey = "PIXELDRAIN_API_KEY"
)

// GlobalFlags manages global flags.
var GlobalFlags = []cli.Flag{
	&cli.StringFlag{
		Name:        FlagAPIKey,
		Usage:       "an API `key` for PixelDrain",
		EnvVars:     []string{EnvPixeldrainApiKey},
		DefaultText: "empty",
	},
}
