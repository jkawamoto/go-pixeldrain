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
	"path"

	"github.com/hashicorp/go-multierror"
	"github.com/urfave/cli/v2"

	"github.com/jkawamoto/go-pixeldrain/cmd/pd/status"
	"github.com/jkawamoto/go-pixeldrain/pkg/pixeldrain"
	"github.com/jkawamoto/go-pixeldrain/pkg/pixeldrain/client"
)

func CmdUpload(c *cli.Context) error {
	if c.NArg() != 1 {
		_, _ = fmt.Printf("expected 1 argument. (%d given)\n", c.NArg())
		return cli.ShowSubcommandHelp(c)
	}

	pd := pixeldrain.New()
	ctx := context.Background()
	if c.Args().First() == "-" {
		id, err := pd.Upload(ctx, os.Stdin, c.String("name"))
		if err != nil {
			return cli.Exit(err, status.APIError)
		}
		if err := printID(id); err != nil {
			return cli.Exit(err, status.IOError)
		}
		return nil
	}

	fp, err := os.Open(c.Args().First())
	if err != nil {
		return cli.Exit(err, status.InvalidArgument)
	}
	defer func() {
		if e := fp.Close(); e != nil && !errors.Is(err, os.ErrClosed) {
			err = multierror.Append(err, e)
		}
	}()

	id, err := pd.Upload(ctx, fp, c.String("name"))
	if err != nil {
		return cli.Exit(err, status.APIError)
	}

	if err := printID(id); err != nil {
		return cli.Exit(err, status.IOError)
	}
	return nil
}

func printID(id string) error {
	_, err := fmt.Printf("https://%v\n", path.Join(client.DefaultHost, client.DefaultBasePath, "file", id))
	return err
}
