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
	"github.com/go-openapi/swag"
	"gopkg.in/cheggaaa/pb.v1"

	"github.com/jkawamoto/go-pixeldrain/pkg/pixeldrain/client/file"
)

// File is the interface the uploading file needs to implement.
type File interface {
	io.Reader
	Name() string
	Stat() (os.FileInfo, error)
}

// Upload the given file to PixelDrain under the given context.
// After the upload succeeds, an ID associated with the uploaded file will be returned.
func (pd *Pixeldrain) Upload(ctx context.Context, f File) (string, error) {
	info, err := f.Stat()
	if err != nil {
		return "", err
	}

	bar := pb.New(int(info.Size())).SetUnits(pb.U_BYTES)
	bar.Output = pd.Stderr
	bar.Start()
	defer bar.Finish()

	name := filepath.Base(f.Name())
	res, err := pd.Client.File.UploadFile(
		file.NewUploadFileParamsWithContext(ctx).
			WithFile(runtime.NamedReader(name, bar.NewProxyReader(f))).
			WithName(swag.String(name)),pd.authInfoWriter,
	)
	if err != nil {
		return "", NewError(err)
	}

	return res.Payload.ID, nil
}
