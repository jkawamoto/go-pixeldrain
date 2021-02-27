// upload.go
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

	"github.com/go-openapi/runtime"
	"gopkg.in/cheggaaa/pb.v1"

	"github.com/jkawamoto/go-pixeldrain/pkg/pixeldrain/client/file"
)

// Upload the given file to PixelDrain under the given context.
// If a name is given, the uploaded file will be renamed.
// After the upload succeeds, an ID associated with the uploaded file will be returned.
func (pd *Pixeldrain) Upload(ctx context.Context, fp *os.File, name string) (id string, err error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if name == "" {
		name = filepath.Base(fp.Name())
	}

	info, err := fp.Stat()
	if err != nil {
		return
	}
	bar := pb.New(int(info.Size())).SetUnits(pb.U_BYTES)
	bar.Output = pd.Stderr
	bar.Start()
	defer bar.Finish()

	return pd.UploadRaw(ctx, bar.NewProxyReader(fp), name)

}

func (pd *Pixeldrain) UploadRaw(ctx context.Context, r io.Reader, name string) (id string, err error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	res, err := pd.Client.File.UploadFile(
		file.NewUploadFileParamsWithContext(ctx).WithFile(runtime.NamedReader(name, r)).WithName(&name))
	if err != nil {
		return
	}

	id = res.Payload.ID
	return

}
