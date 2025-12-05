// auth_test.go
//
// Copyright (c) 2018-2025 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package testutil

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/go-openapi/runtime"
)

func ExpectAuthInfoWritesAPIKey(t *testing.T, authInfo runtime.ClientAuthInfoWriter, apiKey string) {
	t.Helper()

	req := &runtime.TestClientRequest{}
	if err := authInfo.AuthenticateRequest(req, nil); err != nil {
		t.Fatal(err)
	}
	s, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(req.Headers.Get(runtime.HeaderAuthorization), "Basic "))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(s), apiKey) {
		t.Errorf("expect api key, got %v", string(s))
	}
}
