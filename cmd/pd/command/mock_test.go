// mock_test.go
//
// Copyright (c) 2018-2023 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package command

import (
	"testing"

	"github.com/jkawamoto/go-pixeldrain"
	"github.com/jkawamoto/go-pixeldrain/client"
	"github.com/jkawamoto/go-pixeldrain/client/file"
	"github.com/jkawamoto/go-pixeldrain/client/list"
	"github.com/jkawamoto/go-pixeldrain/client/user"
)

// ClientService defines interface that a mock service should implement.
//
//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source=mock_test.go -destination=mock/mock_service.go -package=mock
type ClientService interface {
	file.ClientService
	list.ClientService
	user.ClientService
}

// RegisterMock replaces pixeldrain.Default with a mocked client based on the given mock.
// It will be reset after the given test ends.
func RegisterMock(t *testing.T, mock ClientService) {
	t.Helper()

	old := pixeldrain.Default
	t.Cleanup(func() {
		pixeldrain.Default = old
	})
	pixeldrain.Default = &client.PixeldrainAPI{
		File: mock,
		List: mock,
		User: mock,
	}
}
