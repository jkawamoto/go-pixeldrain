// upload.go
//
// Copyright (c) 2018-2021 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package command

import (
	"errors"
	"fmt"
	"os"

	"github.com/hashicorp/go-multierror"
	"github.com/urfave/cli/v2"

	"github.com/jkawamoto/go-pixeldrain"
	"github.com/jkawamoto/go-pixeldrain/cmd/pd/status"
)

type renamedFile struct {
	*os.File
	name string
}

func (f *renamedFile) Name() string {
	if f.name != "" {
		return f.name
	}
	return f.File.Name()
}

func CmdUpload(c *cli.Context) error {
	if c.NArg() != 1 {
		_, _ = fmt.Printf("expected 1 argument. (%d given)\n", c.NArg())
		return cli.ShowSubcommandHelp(c)
	}

	pd := pixeldrain.New(c.String("api-key"))
	if c.Args().First() == "-" {
		id, err := pd.Upload(c.Context, &renamedFile{File: os.Stdin, name: c.String("name")})
		if err != nil {
			return cli.Exit(err, status.APIError)
		}
		fmt.Println(pd.DownloadURL(id))
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

	id, err := pd.Upload(c.Context, &renamedFile{File: fp, name: c.String("name")})
	if err != nil {
		return cli.Exit(err, status.APIError)
	}

	fmt.Println(pd.DownloadURL(id))
	return nil
}
