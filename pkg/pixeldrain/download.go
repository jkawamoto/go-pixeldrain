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
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-multierror"
	"gopkg.in/cheggaaa/pb.v1"

	"github.com/jkawamoto/go-pixeldrain/pkg/pixeldrain/client/file"
)

func (pd *Pixeldrain) Download(ctx context.Context, url, dir string) error {
	id := url[strings.LastIndex(url, "/")+1:]

	info, err := pd.Client.File.GetFileInfo(file.NewGetFileInfoParamsWithContext(ctx).WithID(id))
	if err != nil {
		var e ErrorResponse
		if errors.As(err, &e) {
			return NewAPIError(e)
		}
		return err
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

	bar := pb.New(int(info.Payload.Size)).SetUnits(pb.U_BYTES)
	bar.Output = pd.Stderr
	bar.Start()
	defer bar.Finish()

	_, err = pd.Client.File.DownloadFile(
		file.NewDownloadFileParamsWithContext(ctx).WithID(info.Payload.ID), nil, io.MultiWriter(out, bar))
	return err
}
