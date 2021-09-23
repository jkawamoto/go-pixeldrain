// doc.go
//
// Copyright (c) 2018-2021 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

// pd is a Pixeldrain client.
//
// Usage
//
// pd has three subcommands: upload, download, and create-list.
//
// upload command uploads a file to Pixeldrain and shows a URL to it.
//   $ pd upload <file path>
// The uploaded file has the same name as the given file. -n and --name options overwrite the file names.
//
// If uploading file is given via STDIN, use - instead of a file path.
// In this case either `-n` or `--name` option is mandatory.
//
// For example, this command reads file1.txt and uploads it with name uploaded.txt:
//   $ cat file1.txt | pd upload --name uploaded.txt -
//
//
// download command downloads a file from Pixeldrain and writes it to STDOUT.
//   $ pd download <file ID | URL>
//
// If -o option is given with a directory path, the downloaded file is stored in
// the directory instead of writing to STDOUT.
//
// For example, this command downloads a file FILE_ID into ~/Download:
//   $ pd download FILE_ID -o ~/Download
//
//
// This application supports uploading a file from STDIN and downloading a file to STDOUT.
// With tar command, it's also able to upload/download directories.
//
// For example, this command uploads ~/Documents directory:
//   $ tar zcf - ~/Documents | pd upload -n documents.tar.gz -
//
// and this command downloads the file:
//   $ pd download <file id> | tar zxf - -C ~/Downloads
//
//
// Installation
//
// If you're a Homebrew or Linuxbrew user, you can install this app by the following commands:
//   $ brew tap jkawamoto/pixeldrain
//   $ brew install pixeldrain
//
// To build the newest version, use go get command:
//    $ go get github.com/jkawamoto/go-pixeldrain
//
// Otherwise, compiled binaries are also available at https://github.com/jkawamoto/go-pixeldrain/releases.
//
//
// License
//
// This software is released under the MIT License.
//
package main
