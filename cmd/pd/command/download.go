// download.go
//
// Copyright (c) 2018-2023 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package command

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/cheggaaa/pb/v3"
	"github.com/go-openapi/swag"
	"github.com/urfave/cli/v2"

	"github.com/jkawamoto/go-pixeldrain"
	"github.com/jkawamoto/go-pixeldrain/client/file"
	"github.com/jkawamoto/go-pixeldrain/cmd/pd/auth"
	"github.com/jkawamoto/go-pixeldrain/cmd/pd/status"
)

func download(ctx *cli.Context, url, dir string) error {
	id := url[strings.LastIndex(url, "/")+1:]

	info, err := pixeldrain.Default.File.GetFileInfo(
		file.NewGetFileInfoParamsWithContext(ctx.Context).WithID(id),
		auth.Extract(ctx.Context),
	)
	if err != nil {
		return pixeldrain.NewError(err)
	}

	out := ctx.App.Writer
	if dir != "" {
		var fp io.WriteCloser
		fp, err = os.OpenFile(filepath.Join(dir, info.Payload.Name), os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer func() {
			err = errors.Join(err, fp.Close())
		}()
		out = fp
	}

	bar := pb.New64(info.Payload.Size)
	bar.Set(pb.SIBytesPrefix, true)
	bar.Set("prefix", info.Payload.Name+" ")
	bar.SetWriter(ctx.App.ErrWriter)
	bar.Start()
	defer bar.Finish()

	_, err = pixeldrain.Default.File.DownloadFile(
		file.NewDownloadFileParamsWithContext(ctx.Context).WithID(swag.StringValue(info.Payload.ID)),
		auth.Extract(ctx.Context),
		bar.NewProxyWriter(out),
	)
	if err != nil {
		return pixeldrain.NewError(err)
	}
	return nil

}

func CmdDownload(c *cli.Context) error {
	if c.NArg() != 1 {
		return cli.Exit(fmt.Sprintf("expected 1 argument but %d given", c.NArg()), status.InvalidArgument)
	}

	url := c.Args().First()
	dir := c.String("output")

	err := download(c, url, dir)
	if err != nil {
		return cli.Exit(err, status.APIError)
	}

	return nil
}
