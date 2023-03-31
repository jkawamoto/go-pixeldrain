// context_test.go
//
// Copyright (c) 2018-2023 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package client

import (
	"context"
	"testing"
)

func TestContext(t *testing.T) {
	cli := New("")

	t.Run("with a client", func(t *testing.T) {
		ctx := ToContext(context.Background(), cli)
		res := Extract(ctx)
		if res != cli {
			t.Errorf("expect %v, got %v", cli, res)
		}
	})

	t.Run("without a client", func(t *testing.T) {
		res := Extract(context.Background())
		if res != nil {
			t.Errorf("expect nil, got %v", res)
		}
	})
}
