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

	"github.com/urfave/cli/v2"

	"github.com/jkawamoto/go-pixeldrain"
	"github.com/jkawamoto/go-pixeldrain/cmd/pd/status"
)

func CmdCreateList(c *cli.Context) error {
	if c.NArg() == 0 {
		_, _ = fmt.Println("expected at least 1 argument.")
		return cli.ShowSubcommandHelp(c)
	}

	pd := pixeldrain.New(c.String("api-key"))
	id, err := pd.CreateList(
		c.Context, c.String("title"), c.String("description"),
		append([]string{c.Args().First()}, c.Args().Tail()...))
	if err != nil {
		return cli.Exit(err, status.APIError)
	}

	if _, err := fmt.Println(pd.ListURL(id)); err != nil {
		return cli.Exit(err, status.IOError)
	}
	return nil
}
