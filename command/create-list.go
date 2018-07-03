// create-list.go
//
// Copyright (c) 2018 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package command

import (
	"context"
	"fmt"

	"github.com/jkawamoto/go-pixeldrain/pixeldrain"
	"github.com/urfave/cli"
)

func CmdCreateList(c *cli.Context) (err error) {

	if c.NArg() == 0 {
		fmt.Println("expected at least 1 argument.")
		return cli.ShowSubcommandHelp(c)
	}

	id, err := pixeldrain.CreateList(
		context.Background(), nil, c.String("title"), c.String("description"),
		append([]string{c.Args().First()}, c.Args().Tail()...))
	if err != nil {
		return cli.NewExitError(err, 2)
	}

	fmt.Println(fmt.Sprintf("https://sia.pixeldrain.com/l/%s", id))
	return

}
