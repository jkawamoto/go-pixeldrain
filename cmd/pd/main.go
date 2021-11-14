// main.go
//
// Copyright (c) 2018-2021 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package main

import (
	"context"
	"os"

	"github.com/urfave/cli/v2"
)

const (
	// Name defines the basename of this program.
	Name = "pd"
	// Version defines current version number.
	Version = "0.5.1"
)

func main() {
	app := cli.NewApp()
	app.Name = Name
	app.Version = Version
	app.Authors = []*cli.Author{
		{
			Name:  "Junpei Kawamoto",
			Email: "kawamoto.junpei@gmail.com",
		},
	}
	app.Usage = "A Pixeldrain client"

	app.Flags = GlobalFlags
	app.Commands = Commands
	app.CommandNotFound = commandNotFound
	app.EnableBashCompletion = true

	_ = app.RunContext(context.Background(), os.Args)
}
