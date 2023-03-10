// doc.go
//
// Copyright (c) 2018-2023 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

// Package pixeldrain provides a Pixeldrain client.
//
// New function takes an API key and returns a new client. The API key can be empty.
// See https://pixeldrain.com/api to find more information about API key.
//
// The client provides Upload and Download functions. See each example for more information.
package pixeldrain

//go:generate go run github.com/go-swagger/go-swagger/cmd/swagger@v0.30.4 generate client -f https://raw.githubusercontent.com/jkawamoto/pixeldrain-swagger/master/swagger.yaml -t .
