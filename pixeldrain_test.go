// pixeldrain_test.go
//
// Copyright (c) 2018-2023 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package pixeldrain

import (
	"path"
	"testing"

	"github.com/jkawamoto/go-pixeldrain/client"
)

func TestDownloadURL(t *testing.T) {
	id := "test-id"
	expect := "https://" + path.Join(client.DefaultHost, client.DefaultBasePath, "file", id)

	if res := DownloadURL(id); res != expect {
		t.Errorf("expect %v, got %v", expect, res)
	}
}

func TestListURL(t *testing.T) {
	id := "test-id"
	expect := "https://" + path.Join(client.DefaultHost, "l", id)

	if res := ListURL(id); res != expect {
		t.Errorf("expect %v, got %v", expect, res)
	}
}
