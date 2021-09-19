// upload.go
//
// Copyright (c) 2018-2021 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package command

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/hashicorp/go-multierror"
	"github.com/urfave/cli"

	"github.com/jkawamoto/go-pixeldrain/cmd/pd/status"
	"github.com/jkawamoto/go-pixeldrain/pkg/pixeldrain"
)

func CmdUpload(c *cli.Context) error {
	if c.NArg() != 1 {
		_, _ = fmt.Printf("expected 1 argument. (%d given)\n", c.NArg())
		return cli.ShowSubcommandHelp(c)
	}

	pd := pixeldrain.New()
	ctx := context.Background()
	if c.Args().First() == "-" {
		id, err := pd.UploadRaw(ctx, os.Stdin, c.String("name"))
		if err != nil {
			return cli.NewExitError(err, status.APIError)
		}
		printID(id)
		return nil
	}

	fp, err := os.Open(c.Args().First())
	if err != nil {
		return cli.NewExitError(err, status.InvalidArgument)
	}
	defer func() {
		if e := fp.Close(); e != nil && !errors.Is(err, os.ErrClosed) {
			err = multierror.Append(err, e)
		}
	}()

	id, err := pd.Upload(ctx, fp, c.String("name"))
	if err != nil {
		return cli.NewExitError(err, status.APIError)
	}

	printID(id)
	return nil
}

func printID(id string) {
	_, _ = fmt.Println(pixeldrain.DownloadEndpoint + "/" + id)
}
