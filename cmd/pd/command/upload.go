// upload.go
//
// Copyright (c) 2018-2024 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package command

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"filippo.io/age"
	"filippo.io/age/agessh"
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

var ErrUnsupportedPublicKey = errors.New("not supported public key")

// upload uploads the given contents with the given name via the given client, and returns the ID of the file.
// If path is "-", read contents from stdin. If name is empty, use the base name of the path.
func upload(ctx *cli.Context, path, name string, recipients []age.Recipient) (_ string, err error) {
	if name == "" {
		name = filepath.Base(path)
	}

	var r io.ReadCloser
	if path == "-" {
		r = io.NopCloser(ctx.App.Reader)
	} else {
		r, err = os.Open(path)
		if err != nil {
			return "", err
		}
		defer func() {
			if e := r.Close(); e != nil && !errors.Is(e, os.ErrClosed) {
				err = errors.Join(err, e)
			}
		}()

		var info os.FileInfo
		info, err = os.Stat(path)
		if err != nil {
			return "", err
		}
		bar := pb.New(int(info.Size()))
		bar.Set(pb.SIBytesPrefix, true)
		bar.Set("prefix", name+" ")
		bar.SetWriter(ctx.App.ErrWriter)
		bar.Start()
		defer bar.Finish()

		r = bar.NewProxyReader(r)
	}

	if len(recipients) != 0 {
		r = Encrypt(r, recipients)
		defer func() {
			err = errors.Join(err, r.Close())
		}()
		name = name + AgeExt
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

	recipients, err := parseRecipients(c.StringSlice(FlagRecipient))
	if err != nil {
		return cli.Exit(err, status.InvalidArgument)
	}
	if name := c.String(FlagRecipientFile); name != "" {
		rs, err := parseRecipientFile(name)
		if err != nil {
			return cli.Exit(err, status.InvalidArgument)
		}
		recipients = append(recipients, rs...)
	}

	var ids []string
	for i := 0; i != c.NArg(); i++ {
		path, name, err := parseArgument(c.Args().Get(i))
		if err != nil {
			return cli.Exit(
				fmt.Errorf("failed to parse argument %q: %w", c.Args().Get(i), err), status.InvalidArgument)
		}
		id, err := upload(c, path, name, recipients)
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

func parseRecipients(recipients []string) ([]age.Recipient, error) {
	res := make([]age.Recipient, len(recipients))
	for i, s := range recipients {
		switch {
		case strings.HasPrefix(s, "age1"):
			r, err := age.ParseX25519Recipient(s)
			if err != nil {
				return nil, err
			}
			res[i] = r
		case strings.HasPrefix(s, "ssh-"):
			r, err := agessh.ParseRecipient(s)
			if err != nil {
				return nil, err
			}
			res[i] = r
		default:
			return nil, ErrUnsupportedPublicKey
		}
	}
	return res, nil
}

func parseRecipientFile(name string) (_ []age.Recipient, err error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = errors.Join(err, f.Close())
	}()

	return age.ParseRecipients(f)
}

func parseArgument(arg string) (path string, name string, _ error) {
	if strings.HasPrefix(arg, "-:") {
		name = strings.TrimPrefix(arg, "-:")
		if strings.HasPrefix(name, "\"") && strings.HasSuffix(name, "\"") {
			name = name[1 : len(name)-1]
		}
		return "-", name, nil
	}

	volume := filepath.VolumeName(arg)
	arg = strings.TrimPrefix(arg, volume)

	r := csv.NewReader(strings.NewReader(arg))
	r.Comma = ':'
	res, err := r.Read()
	if err != nil {
		return "", "", err
	}

	switch len(res) {
	case 0:
		return "", "", nil
	case 1:
		return volume + res[0], "", nil
	default:
		return volume + res[0], res[1], nil
	}
}
