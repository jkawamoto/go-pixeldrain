// list.go
//
// Copyright (c) 2018-2023 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package command

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/jkawamoto/go-pixeldrain/cmd/client"
	"github.com/jkawamoto/go-pixeldrain/cmd/pd/status"
)

func CmdCreateList(c *cli.Context) error {
	if c.NArg() == 0 {
		_, _ = fmt.Println("expected at least 1 argument.")
		return cli.ShowSubcommandHelp(c)
	}

	pd := client.Extract(c.Context)
	id, err := pd.CreateList(
		c.Context, c.String("title"), append([]string{c.Args().First()}, c.Args().Tail()...))
	if err != nil {
		return cli.Exit(err, status.APIError)
	}

	if _, err := fmt.Println(pd.ListURL(id)); err != nil {
		return cli.Exit(err, status.IOError)
	}
	return nil
}
