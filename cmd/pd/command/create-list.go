// create-list.go
//
// Copyright (c) 2018-2021 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package command

import (
	"context"
	"fmt"
	"path"

	"github.com/urfave/cli"

	"github.com/jkawamoto/go-pixeldrain/pkg/pixeldrain"
	"github.com/jkawamoto/go-pixeldrain/pkg/pixeldrain/client"
)

func CmdCreateList(c *cli.Context) (err error) {

	if c.NArg() == 0 {
		fmt.Println("expected at least 1 argument.")
		return cli.ShowSubcommandHelp(c)
	}

	id, err := pixeldrain.New().CreateList(
		context.Background(), c.String("title"), c.String("description"),
		append([]string{c.Args().First()}, c.Args().Tail()...))
	if err != nil {
		return cli.NewExitError(err, 2)
	}

	fmt.Println(fmt.Sprintf("https://%v", path.Join(client.DefaultHost, "l", id)))
	return

}
