// example_test.go
//
// Copyright (c) 2018-2021 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package pixeldrain_test

import (
	"context"
	"fmt"
	"os"

	"github.com/jkawamoto/go-pixeldrain"
)

// Since os.File implements File interface, you can open the file and pass it to upload a file.
func ExamplePixeldrain_Upload() {
	// Open the target file.
	f, err := os.Open("example_test.go")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Create a client and upload the file. API key can be empty.
	cli := pixeldrain.New("YOUR API KEY IF NECESSARY")
	id, err := cli.Upload(context.Background(), f)
	if err != nil {
		panic(err)
	}

	// The file ID is used to download the file with Download function.
	fmt.Println("File ID:", id)
	// DownloadURL returns the file URL so that users can open it with their browsers, wget, curl, etc.
	fmt.Println("Download URL:", cli.DownloadURL(id))
}

// Download has two behaviours. If dir is given, the downloaded file will be stored in the given directory.
// Otherwise, it will be written in a stream.
func ExamplePixeldrain_Download() {
	// Create a client and upload the file. API key can be empty.
	cli := pixeldrain.New("YOUR API KEY IF NECESSARY")

	// url can be a URL or just a file ID.
	url := "https://pixeldrain.com/api/file/FILE_ID" // or "FILE_ID"

	// If a directory path is given, the downloaded file will be stored in the directory.
	if err := cli.Download(context.Background(), url, "~/Downloads"); err != nil {
		panic(err)
	}

	// If no directory path is given, the downloaded file will be written in cli.Stdout i.e. STDOUT by default.
	if err := cli.Download(context.Background(), url, ""); err != nil {
		panic(err)
	}
}
