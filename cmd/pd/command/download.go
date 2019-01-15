/*
 * download.go
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
	"github.com/jkawamoto/go-pixeldrain/pixeldrain"
	"github.com/urfave/cli"
)

func CmdDownload(c *cli.Context) (err error) {

	if c.NArg() != 1 {
		fmt.Println(fmt.Sprintf("expected 1 argument. (%d given)", c.NArg()))
		return cli.ShowSubcommandHelp(c)
	}

	url := c.Args().First()
	dir := c.String("output")

	err = pixeldrain.New().Download(context.Background(), url, dir)
	if err != nil {
		return cli.NewExitError(err, 2)
	}

	return

}
