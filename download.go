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

	"github.com/cheggaaa/pb/v3"
	"github.com/go-openapi/swag"

	"github.com/jkawamoto/go-pixeldrain/client/file"
)

// Download the file associated with the given url or file ID. If dir is given, the downloaded file is stored into
// the directory. Otherwise, it is written in pd.Stdout.
func (pd *Pixeldrain) Download(ctx context.Context, url, dir string) (err error) {
	id := url[strings.LastIndex(url, "/")+1:]

	info, err := pd.cli.File.GetFileInfo(file.NewGetFileInfoParamsWithContext(ctx).WithID(id), pd.authInfoWriter)
	if err != nil {
		return NewError(err)
	}

	out := pd.Stdout
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
	bar.SetWriter(pd.Stderr)
	bar.Start()
	defer bar.Finish()

	_, err = pd.cli.File.DownloadFile(
		file.NewDownloadFileParamsWithContext(ctx).WithID(swag.StringValue(info.Payload.ID)),
		pd.authInfoWriter,
		bar.NewProxyWriter(out),
	)
	if err != nil {
		return NewError(err)
	}
	return nil
}
