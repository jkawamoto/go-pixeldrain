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

func TestIsDownloadURL(t *testing.T) {
	cases := []struct {
		url    string
		expect bool
	}{
		{
			url:    "https://" + path.Join(client.DefaultHost, client.DefaultBasePath, "file", "123"),
			expect: true,
		},
		{
			url:    "https://" + path.Join(client.DefaultHost, client.DefaultBasePath, "l", "123"),
			expect: false,
		},
		{
			url:    "https://" + path.Join(client.DefaultHost, "file", "123"),
			expect: false,
		},
		{
			url:    "https://" + path.Join(client.DefaultHost, "l", "123"),
			expect: false,
		},
	}
	for _, c := range cases {
		t.Run(c.url, func(t *testing.T) {
			res, err := IsDownloadURL(c.url)
			if err != nil {
				t.Fatal(err)
			}
			if res != c.expect {
				t.Errorf("expect %v, got %v", c.expect, res)
			}
		})
	}
}

func TestIsListURL(t *testing.T) {
	cases := []struct {
		url    string
		expect bool
	}{
		{
			url:    "https://" + path.Join(client.DefaultHost, client.DefaultBasePath, "file", "123"),
			expect: false,
		},
		{
			url:    "https://" + path.Join(client.DefaultHost, client.DefaultBasePath, "l", "123"),
			expect: false,
		},
		{
			url:    "https://" + path.Join(client.DefaultHost, "file", "123"),
			expect: false,
		},
		{
			url:    "https://" + path.Join(client.DefaultHost, "l", "123"),
			expect: true,
		},
	}
	for _, c := range cases {
		t.Run(c.url, func(t *testing.T) {
			res, err := IsListURL(c.url)
			if err != nil {
				t.Fatal(err)
			}
			if res != c.expect {
				t.Errorf("expect %v, got %v", c.expect, res)
			}
		})
	}
}
