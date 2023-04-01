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

	"github.com/urfave/cli/v2"

	"github.com/jkawamoto/go-pixeldrain"
	"github.com/jkawamoto/go-pixeldrain/cmd/client"
)

type file interface {
	pixeldrain.File
	io.Closer
}

type fileWrapper struct {
	io.ReadCloser
}

// Name returns the name of this file.
func (w fileWrapper) Name() string {
	if namer, ok := w.ReadCloser.(interface {
		Name() string
	}); ok {
		return namer.Name()
	}
	return ""
}

// Stat returns os.FileInfo of this file.
func (w fileWrapper) Stat() (os.FileInfo, error) {
	if stater, ok := w.ReadCloser.(interface {
		State() (os.FileInfo, error)
	}); ok {
		return stater.State()
	}
	return nil, fmt.Errorf("couldn't find file info")
}

var _ file = (*fileWrapper)(nil)

type renamedFile struct {
	file
	name string
}

func (f *renamedFile) Name() string {
	return f.name
}

// upload uploads the given contents with the given name via the given client, and returns the ID of the file.
// If path is "-", read contents from stdin. If name is empty, use the base name of the path.
func upload(ctx *cli.Context, cli client.Client, path, name string) (_ string, err error) {
	if name == "" {
		name = filepath.Base(path)
	}

	var f file
	if path == "-" {
		f = fileWrapper{ReadCloser: io.NopCloser(ctx.App.Reader)}
	} else {
		f, err = os.Open(path)
		defer func() {
			e := f.Close()
			if e != nil && !errors.Is(e, os.ErrClosed) {
				err = errors.Join(err, e)
			}
		}()
	}

	return cli.Upload(ctx.Context, &renamedFile{
		file: f,
		name: name,
	})
}

func CmdUpload(c *cli.Context) error {
	if c.NArg() == 0 {
		return fmt.Errorf("expected 1 argument but %d given: %w", c.NArg(), ErrNotEnoughArguments)
	}

	pd := client.Extract(c.Context)

	var ids []string
	for i := 0; i != c.NArg(); i++ {
		path, name, _ := strings.Cut(c.Args().Get(i), ":")

		id, err := upload(c, pd, path, name)
		if err != nil {
			return err
		}

		ids = append(ids, id)
	}

	if len(ids) == 1 {
		_, _ = fmt.Fprintln(c.App.Writer, pd.DownloadURL(ids[0]))
		return nil
	}

	album := c.String(FlagAlbumName)
	if album == "" {
		album = fmt.Sprintf("album-%x", time.Now().Unix())
	}
	id, err := pd.CreateList(c.Context, album, ids)
	if err != nil {
		return err
	}

	_, _ = fmt.Fprintln(c.App.Writer, pd.ListURL(id))
	return nil
}
