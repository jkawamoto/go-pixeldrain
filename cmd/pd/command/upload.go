/*
 * upload.go
 *
 * Copyright (c) 2018-2019 Junpei Kawamoto
 *
 * This software is released under the MIT License.
 *
 * http://opensource.org/licenses/mit-license.php
 */

package command

import (
	"context"
	"fmt"
	"os"

	"github.com/jkawamoto/go-pixeldrain/pixeldrain"
	"github.com/urfave/cli"
)

func CmdUpload(c *cli.Context) (err error) {

	if c.NArg() != 1 {
		fmt.Println(fmt.Sprintf("expected 1 argument. (%d given)", c.NArg()))
		return cli.ShowSubcommandHelp(c)
	}

	fp, err := os.Open(c.Args().First())
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	defer fp.Close()

	id, err := pixeldrain.Upload(context.Background(), nil, fp, c.String("name"))
	if err != nil {
		return cli.NewExitError(err, 2)
	}

	fmt.Println(fmt.Sprintf("https://sia.pixeldrain.com/api/file/%s", id))
	return

}
