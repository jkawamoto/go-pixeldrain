// context_test.go
//
// Copyright (c) 2018-2025 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package auth

import (
	"context"
	"reflect"
	"runtime"
	"testing"

	"github.com/go-openapi/runtime/client"
)

func funcName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func TestContext(t *testing.T) {
	info := client.BasicAuth("", "")

	t.Run("with auth info", func(t *testing.T) {
		ctx := ToContext(context.Background(), info)
		res := Extract(ctx)
		if funcName(res) != funcName(info) {
			t.Errorf("expect %v, got %v", info, res)
		}
	})

	t.Run("without auth info", func(t *testing.T) {
		res := Extract(context.Background())
		if res != nil {
			t.Errorf("expect nil, got %v", res)
		}
	})
}
