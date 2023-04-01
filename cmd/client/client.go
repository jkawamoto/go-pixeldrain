// client.go
//
// Copyright (c) 2018-2023 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package client

import (
	"context"

	"github.com/jkawamoto/go-pixeldrain"
)

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source=client.go -destination=mock/mock_client.go -package=mock

type Client interface {
	Upload(ctx context.Context, f pixeldrain.File) (string, error)
	Download(ctx context.Context, url, dir string) error
	CreateList(ctx context.Context, title string, items []string) (string, error)
	DownloadURL(id string) string
	ListURL(id string) string
}

func New(apiKey string) Client {
	return pixeldrain.New(apiKey)
}
