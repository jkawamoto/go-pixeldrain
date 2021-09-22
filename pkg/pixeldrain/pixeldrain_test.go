// pixeldrain_test.go
//
// Copyright (c) 2018-2021 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package pixeldrain

import (
	"encoding/base64"
	"os"
	"path"
	"testing"

	"github.com/go-openapi/runtime"

	"github.com/jkawamoto/go-pixeldrain/pkg/pixeldrain/client"
)

func authorization(apiKey string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(":"+apiKey))
}

func TestNew(t *testing.T) {
	apiKey := "123abc456"

	cases := []struct {
		name   string
		apiKey string
		verify func(*testing.T, *Pixeldrain)
	}{
		{
			name:   "without api key",
			apiKey: apiKey,
			verify: func(t *testing.T, pd *Pixeldrain) {
				if pd.authInfoWriter == nil {
					t.Fatal("expect authInfoWriter is not nil")
				}

				req := &runtime.TestClientRequest{}
				err := pd.authInfoWriter.AuthenticateRequest(req, nil)
				if err != nil {
					t.Error(err)
				}

				expect := authorization(apiKey)
				if h := req.Headers.Get("Authorization"); h != expect {
					t.Errorf("expect %v, got %v", expect, h)
				}
			},
		},
		{
			name: "without api key",
			verify: func(t *testing.T, pd *Pixeldrain) {
				if pd.authInfoWriter != nil {
					t.Errorf("expect nil")
				}
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			pd := New(c.apiKey)
			if pd.cli != client.Default {
				t.Errorf("Client is %v, expects %v", pd.cli, client.Default)
			}
			if pd.Stderr != os.Stderr {
				t.Errorf("Stderr is %v, expects %v", pd.Stderr, os.Stderr)
			}
			c.verify(t, pd)
		})
	}
}

func TestPixeldrain_DownloadURL(t *testing.T) {
	id := "test-id"
	expect := "https://" + path.Join(client.DefaultHost, client.DefaultBasePath, "file", id)

	pd := &Pixeldrain{}
	if res := pd.DownloadURL(id); res != expect {
		t.Errorf("expect %v, got %v", expect, res)
	}
}
