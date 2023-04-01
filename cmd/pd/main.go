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
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/jkawamoto/go-pixeldrain/cmd/client"
	"github.com/jkawamoto/go-pixeldrain/cmd/pd/command"
	"github.com/jkawamoto/go-pixeldrain/cmd/pd/status"
)

const (
	// Name defines the basename of this program.
	Name = "pd"
	// Version defines current version number.
	Version = "0.5.2"
)

// commandNotFound shows error message and exit when a given command is not found.
func commandNotFound(c *cli.Context, command string) {
	_, _ = fmt.Fprintf(c.App.ErrWriter, "'%s' is not a %s command..\n", command, c.App.Name)
	//_ = cli.ShowAppHelp(c)
	os.Exit(status.InvalidCommand)
}

func initApp() *cli.App {
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
	app.Commands = command.Commands
	app.CommandNotFound = commandNotFound
	app.EnableBashCompletion = true
	app.Before = func(c *cli.Context) error {
		c.Context = client.ToContext(c.Context, client.New(c.String(FlagAPIKey)))
		return nil
	}

	return app
}

func main() {
	app := initApp()

	err := app.RunContext(context.Background(), os.Args)
	if err != nil {
		_, _ = fmt.Fprintf(app.ErrWriter, "failed to run: %v\n", err)
		os.Exit(status.InvalidArgument)
	}
}
