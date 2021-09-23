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

	"github.com/cheggaaa/pb/v3"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"

	"github.com/jkawamoto/go-pixeldrain/client/file"
)

// File is the interface the uploading file needs to implement. os.File implements this interface.
type File interface {
	io.Reader
	// Name returns the name of this file.
	Name() string
	// Stat returns os.FileInfo of this file.
	Stat() (os.FileInfo, error)
}

// Upload the given file to PixelDrain under the given context.
// After the upload succeeds, an ID associated with the uploaded file will be returned.
func (pd *Pixeldrain) Upload(ctx context.Context, f File) (string, error) {
	name := filepath.Base(f.Name())
	info, err := f.Stat()
	if err != nil {
		return "", err
	}

	bar := pb.New(int(info.Size()))
	bar.Set(pb.SIBytesPrefix, true)
	bar.Set("prefix", name+" ")
	bar.SetWriter(pd.Stderr)
	bar.Start()
	defer bar.Finish()

	res, err := pd.cli.File.UploadFile(
		file.NewUploadFileParamsWithContext(ctx).
			WithFile(runtime.NamedReader(name, bar.NewProxyReader(f))).
			WithName(swag.String(name)),
		pd.authInfoWriter,
	)
	if err != nil {
		return "", NewError(err)
	}

	return res.Payload.ID, nil
}
