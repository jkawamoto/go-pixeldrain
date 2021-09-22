// create-list.go
//
// Copyright (c) 2018-2021 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package command

import (
	"fmt"
	"path"

	"github.com/urfave/cli/v2"

	"github.com/jkawamoto/go-pixeldrain/cmd/pd/status"
	"github.com/jkawamoto/go-pixeldrain/pkg/pixeldrain"
	"github.com/jkawamoto/go-pixeldrain/pkg/pixeldrain/client"
)

func CmdCreateList(c *cli.Context) error {
	if c.NArg() == 0 {
		_, _ = fmt.Println("expected at least 1 argument.")
		return cli.ShowSubcommandHelp(c)
	}

	id, err := pixeldrain.New(c.String("api-key")).CreateList(
		c.Context, c.String("title"), c.String("description"),
		append([]string{c.Args().First()}, c.Args().Tail()...))
	if err != nil {
		return cli.Exit(err, status.APIError)
	}

	if _, err := fmt.Printf("https://%v\n", path.Join(client.DefaultHost, "l", id)); err != nil {
		return cli.Exit(err, status.IOError)
	}
	return nil
}
