// download.go
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

	"github.com/urfave/cli/v2"

	"github.com/jkawamoto/go-pixeldrain/cmd/pd/status"
	"github.com/jkawamoto/go-pixeldrain/pkg/pixeldrain"
)

func CmdDownload(c *cli.Context) error {
	if c.NArg() != 1 {
		_, _ = fmt.Printf("expected 1 argument. (%d given)\n", c.NArg())
		return cli.ShowSubcommandHelp(c)
	}

	url := c.Args().First()
	dir := c.String("output")

	err := pixeldrain.New().Download(context.Background(), url, dir)
	if err != nil {
		return cli.Exit(err, status.APIError)
	}

	return nil
}
