// download.go
//
// Copyright (c) 2018-2021 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package pixeldrain

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-multierror"
	"gopkg.in/cheggaaa/pb.v1"

	"github.com/jkawamoto/go-pixeldrain/client/file"
)

// Download the file associated with the given url or file ID. If dir is given, the downloaded file is stored into
// the directory. Otherwise, it is written in pd.Stdout.
func (pd *Pixeldrain) Download(ctx context.Context, url, dir string) error {
	id := url[strings.LastIndex(url, "/")+1:]

	info, err := pd.cli.File.GetFileInfo(file.NewGetFileInfoParamsWithContext(ctx).WithID(id), pd.authInfoWriter)
	if err != nil {
		return NewError(err)
	}

	out := pd.Stdout
	if dir != "" {
		fp, err := os.OpenFile(filepath.Join(dir, info.Payload.Name), os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer func() {
			if e := fp.Close(); e != nil {
				err = multierror.Append(err, e)
			}
		}()
		out = fp
	}

	bar := pb.New(int(info.Payload.Size)).SetUnits(pb.U_BYTES).Prefix(info.Payload.Name)
	bar.Output = pd.Stderr
	bar.Start()
	defer bar.Finish()

	_, err = pd.cli.File.DownloadFile(
		file.NewDownloadFileParamsWithContext(ctx).WithID(info.Payload.ID), pd.authInfoWriter, io.MultiWriter(out, bar))
	if err != nil {
		return NewError(err)
	}
	return nil
}
