// example_test.go
//
// Copyright (c) 2018-2025 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package pixeldrain_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/go-openapi/runtime/client"
	"github.com/go-openapi/swag"

	"github.com/jkawamoto/go-pixeldrain"
	"github.com/jkawamoto/go-pixeldrain/client/file"
)

func Example_upload() {
	// Open the target file.
	f, err := os.Open("example_test.go")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// Upload the file with the default client. API key can be empty.
	res, err := pixeldrain.Default.File.UploadFile(
		file.NewUploadFileParamsWithContext(context.Background()).WithFile(f),
		client.BasicAuth("", "YOUR API KEY IF NECESSARY"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// File ID is used to download the file by this client.
	fmt.Println("File ID:", swag.StringValue(res.Payload.ID))

	// Download URL is for browsers, wget, curl, etc.
	fmt.Println("Download URL:", pixeldrain.DownloadURL(swag.StringValue(res.Payload.ID)))
}

func Example_download() {
	ctx := context.Background()
	id := "FILE_ID"
	apiKey := "YOUR API KEY IF NECESSARY"

	//  Get the file information.
	info, err := pixeldrain.Default.File.GetFileInfo(
		file.NewGetFileInfoParamsWithContext(ctx).WithID(id),
		client.BasicAuth("", apiKey),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Open a file to store the downloaded contents.
	f, err := os.OpenFile(filepath.Join("~/Downloads", info.Payload.Name), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// If a directory path is given, the downloaded file will be stored in the directory.
	_, _, err = pixeldrain.Default.File.DownloadFile(
		file.NewDownloadFileParamsWithContext(ctx).WithID(id),
		client.BasicAuth("", apiKey),
		f,
	)
	if err != nil {
		log.Fatal(err)
	}
}
