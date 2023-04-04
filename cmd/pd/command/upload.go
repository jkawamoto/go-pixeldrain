// upload.go
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
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"
	"github.com/urfave/cli/v2"

	"github.com/jkawamoto/go-pixeldrain"
	"github.com/jkawamoto/go-pixeldrain/client/file"
	"github.com/jkawamoto/go-pixeldrain/client/list"
	"github.com/jkawamoto/go-pixeldrain/cmd/pd/auth"
	"github.com/jkawamoto/go-pixeldrain/cmd/pd/status"
	"github.com/jkawamoto/go-pixeldrain/models"
)

// upload uploads the given contents with the given name via the given client, and returns the ID of the file.
// If path is "-", read contents from stdin. If name is empty, use the base name of the path.
func upload(ctx *cli.Context, path, name string) (_ string, err error) {
	if name == "" {
		name = filepath.Base(path)
	}

	var r io.Reader
	if path == "-" {
		r = ctx.App.Reader
	} else {
		f, e := os.Open(path)
		if e != nil {
			return "", e
		}
		defer func() {
			if e := f.Close(); e != nil && !errors.Is(e, os.ErrClosed) {
				err = errors.Join(err, e)
			}
		}()
		r = f

		info, e := f.Stat()
		if e != nil {
			return "", e
		}
		bar := pb.New(int(info.Size()))
		bar.Set(pb.SIBytesPrefix, true)
		bar.Set("prefix", name+" ")
		bar.SetWriter(ctx.App.ErrWriter)
		bar.Start()
		defer bar.Finish()
	}

	res, err := pixeldrain.Default.File.UploadFile(
		file.NewUploadFileParamsWithContext(ctx.Context).WithFile(runtime.NamedReader(name, r)),
		auth.Extract(ctx.Context),
	)
	if err != nil {
		return "", pixeldrain.NewError(err)
	}

	return swag.StringValue(res.Payload.ID), nil
}

func CmdUpload(c *cli.Context) error {
	if c.NArg() == 0 {
		return cli.Exit(fmt.Sprintf("expected at least 1 argument but %d given", c.NArg()), status.InvalidArgument)
	}

	var ids []string
	for i := 0; i != c.NArg(); i++ {
		path, name, _ := strings.Cut(c.Args().Get(i), ":")

		id, err := upload(c, path, name)
		if err != nil {
			return cli.Exit(fmt.Errorf("failed to upload %v: %w", path, err), status.APIError)
		}

		ids = append(ids, id)
	}

	if len(ids) == 1 {
		_, _ = fmt.Fprintln(c.App.Writer, pixeldrain.DownloadURL(ids[0]))
		return nil
	}

	album := c.String(FlagAlbumName)
	if album == "" {
		album = fmt.Sprintf("album-%x", time.Now().Unix())
	}

	files := make([]*models.ListItem, len(ids))
	for i, id := range ids {
		files[i] = &models.ListItem{
			ID: swag.String(id),
		}
	}
	res, err := pixeldrain.Default.List.CreateFileList(
		list.NewCreateFileListParamsWithContext(c.Context).WithList(&models.CreateFileListRequest{
			Files: files,
			Title: swag.String(album),
		}),
		auth.Extract(c.Context),
	)
	if err != nil {
		return cli.Exit(fmt.Errorf("failed to create an album: %w", pixeldrain.NewError(err)), status.APIError)
	}

	_, _ = fmt.Fprintln(c.App.Writer, pixeldrain.ListURL(swag.StringValue(res.Payload.ID)))
	return nil
}
