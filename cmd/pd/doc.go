// doc.go
//
// Copyright (c) 2018-2021 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

// pd is a Pixeldrain client.
//
// # Usage
//
// pd has two subcommands: upload and download.
//
// upload command uploads files specified by the given paths to Pixeldrain and shows a URL to download them.
//
//	$ pd upload <path[:name]>...
//
// Each path can have an optional name. If a name is given, uploaded file will be renamed with it.
// For example, this command reads img.png and uploads it as another.png:
//
//	$ pd upload img.png:another.png
//
// If path is "-", the uploading file is read from stdin. In this case, it's recommended to give a file name.
// For example, this command reads data from stdin and uploads it as output.log:
//
//	$ pd upload -:output.log
//
// If multiple files are given, an album consists of them will be created. By default, the album has a random name.
// --album flag can specify the name.
// For example, this command uploads two files and creates an album named "screenshots":
//
//	$ pd upload --album screenshots img1.png img2.png
//
// download command downloads a file from Pixeldrain and writes it to STDOUT.
//
//	$ pd download <file ID | URL>
//
// If -o option is given with a directory path, the downloaded file is stored in
// the directory instead of writing to STDOUT.
//
// For example, this command downloads a file FILE_ID into ~/Download:
//
//	$ pd download FILE_ID -o ~/Download
//
// Sine this command supports uploading a file from STDIN and downloading a file to stdin,
// it's also able to upload/download directories with tar command
//
// For example, this command uploads ~/Documents directory:
//
//	$ tar zcf - ~/Documents | pd upload -:documents.tar.gz
//
// and this command downloads the file:
//
//	$ pd download <file id> | tar zxf - -C ~/Downloads
//
// # Installation
//
// If you're a Homebrew or Linuxbrew user, you can install this app by the following commands:
//
//	$ brew tap jkawamoto/pixeldrain
//	$ brew install pixeldrain
//
// To build the newest version, use go get command:
//
//	$ go get github.com/jkawamoto/go-pixeldrain
//
// Otherwise, compiled binaries are also available at https://github.com/jkawamoto/go-pixeldrain/releases.
//
// # License
//
// This software is released under the MIT License.
package main
