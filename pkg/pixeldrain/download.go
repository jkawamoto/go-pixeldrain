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
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/context/ctxhttp"
	"gopkg.in/cheggaaa/pb.v1"

	"github.com/jkawamoto/go-pixeldrain/pkg/pixeldrain/client/file"
)

func (pd *Pixeldrain) Download(ctx context.Context, url, dir string) (err error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	id := url[strings.LastIndex(url, "/")+1:]

	info, err := pd.Client.File.GetFileInfo(file.NewGetFileInfoParamsWithContext(ctx).WithID(id))
	if err != nil {
		return
	}

	fp := pd.Stdout
	if dir != "" {
		fp, err = os.OpenFile(filepath.Join(dir, info.Payload.Name), os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return
		}
		defer func() {
			//noinspection SpellCheckingInspection
			cerr := fp.Close()
			if cerr != nil {
				err = fmt.Errorf("failed to close: %v, the original error was %v", cerr, err)
			}
		}()
	}

	bar := pb.New(int(info.Payload.Size)).SetUnits(pb.U_BYTES)
	bar.Output = pd.Stderr
	bar.Start()
	defer bar.Finish()

	res, err := ctxhttp.Get(ctx, nil, fmt.Sprint(pd.downloadEndpoint, info.Payload.ID))
	if err != nil {
		return
	}
	defer func() {
		//noinspection SpellCheckingInspection
		cerr := res.Body.Close()
		if cerr != nil {
			err = fmt.Errorf("failed to close: %v, the original error was %v", cerr, err)
		}
	}()

	_, err = io.Copy(io.MultiWriter(fp, bar), res.Body)
	return

}
