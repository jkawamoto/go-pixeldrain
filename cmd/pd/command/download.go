// download.go
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

func CmdDownload(c *cli.Context) error {
	if c.NArg() != 1 {
		_, _ = fmt.Printf("expected 1 argument. (%d given)\n", c.NArg())
		return cli.ShowSubcommandHelp(c)
	}

	url := c.Args().First()
	dir := c.String("output")

	pd := client.Extract(c.Context)
	err := pd.Download(c.Context, url, dir)
	if err != nil {
		return cli.Exit(err, status.APIError)
	}

	return nil
}
