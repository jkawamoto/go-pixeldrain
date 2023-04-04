// pixeldrain_api_client_test.go
//
// Copyright (c) 2018-2023 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package client

import "testing"

func TestDefaultSchemes(t *testing.T) {
	if len(DefaultSchemes) == 0 {
		t.Error("expect at least one scheme")
	}
}
