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
	"path"

	"github.com/jkawamoto/go-pixeldrain/client"
	"github.com/jkawamoto/go-pixeldrain/pixeldrain"
	"github.com/urfave/cli"
)

func CmdUpload(c *cli.Context) (err error) {

	if c.NArg() != 1 {
		fmt.Println(fmt.Sprintf("expected 1 argument. (%d given)", c.NArg()))
		return cli.ShowSubcommandHelp(c)
	}

	var id string
	pd := pixeldrain.New()
	ctx := context.Background()
	if c.Args().First() == "-" {

		id, err = pd.UploadRaw(ctx, os.Stdin, c.String("name"))

	} else {

		fp, err := os.Open(c.Args().First())
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		defer func() {
			//noinspection SpellCheckingInspection
			cerr := fp.Close()
			if cerr != nil {
				err = fmt.Errorf("failed to close: %v, the original error was %v", cerr, err)
			}
		}()

		id, err = pd.Upload(ctx, fp, c.String("name"))
		if err != nil {
			return cli.NewExitError(err, 2)
		}

	}

	fmt.Println(fmt.Sprintf("https://%v", path.Join(client.DefaultHost, client.DefaultBasePath, "file", id)))
	return

}
