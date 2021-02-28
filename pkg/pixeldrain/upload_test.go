// upload_test.go
//
// Copyright (c) 2018-2021 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package pixeldrain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"testing"

	"github.com/jkawamoto/go-pixeldrain/pkg/pixeldrain/client"
	"github.com/jkawamoto/go-pixeldrain/pkg/pixeldrain/client/file"
)

const (
	TestFileName = "upload.go"
)

type mockHandler struct {
	expected []byte
	id       string
	name     string
}

func newMockHandler(t *testing.T, file, name, id string) *mockHandler {
	t.Helper()

	fp, err := os.Open(file)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := fp.Close(); err != nil {
			t.Error(err)
		}
	}()

	data, err := ioutil.ReadAll(fp)
	if err != nil {
		t.Fatal(err)
	}

	return &mockHandler{
		expected: data,
		id:       id,
		name:     name,
	}
}

func (m *mockHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/file" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if v := req.FormValue("name"); v != m.name {
		res.WriteHeader(http.StatusBadRequest)
		if _, err := fmt.Fprintf(res, "given file name is %v, want %v", v, m.name); err != nil {
			panic(err)
		}
		return
	}

	fp, _, err := req.FormFile("file")
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		if _, err := fmt.Fprintln(res, "cannot get the file:", err.Error()); err != nil {
			panic(err)
		}
		return
	}
	defer func() {
		if err := fp.Close(); err != nil {
			panic(err)
		}
	}()

	data, err := ioutil.ReadAll(fp)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		if _, err := fmt.Fprintln(res, "Cannot read the uploaded file:", err.Error()); err != nil {
			panic(err)
		}
		return
	}

	if !reflect.DeepEqual(data, m.expected) {
		res.WriteHeader(http.StatusInternalServerError)
		if _, err := fmt.Fprintln(res, "Uploaded file is broken"); err != nil {
			panic(err)
		}
		return
	}

	res.Header().Set(ContentType, "application/json")
	res.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(res).Encode(&file.UploadFileCreatedBody{
		ID:      m.id,
		Success: true,
	}); err != nil {
		panic(err)
	}
}

func TestUpload(t *testing.T) {
	id := "test-id"
	cases := []struct {
		file   string
		rename string
		expect string
	}{
		{file: TestFileName, rename: "", expect: TestFileName},
		{file: TestFileName, rename: "another-expect", expect: "another-expect"},
		{file: "../../cmd/pd/command/upload.go", rename: "", expect: TestFileName},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("%+v", c), func(t *testing.T) {
			server := httptest.NewServer(newMockHandler(t, c.file, c.expect, id))
			defer server.Close()

			pd := New()
			u, err := url.Parse(server.URL)
			if err != nil {
				t.Fatal("Cannot parse a URL:", err)
			}
			pd.Client = client.NewHTTPClientWithConfig(nil, &client.TransportConfig{
				Host:     u.Host,
				BasePath: "/",
				Schemes:  []string{"http"},
			})

			fp, err := os.Open(c.file)
			if err != nil {
				t.Fatal("Failed to open the file:", err)
			}
			defer func() {
				if err := fp.Close(); err != nil && !errors.Is(err, os.ErrClosed) {
					t.Error(err)
				}
			}()

			res, err := pd.Upload(context.Background(), fp, c.rename)
			if err != nil {
				t.Fatal("failed to upload a file:", err.Error())
			}
			if res != id {
				t.Errorf("received ID = %v, want %v", res, id)
			}
		})
	}
}

func TestUploadRaw(t *testing.T) {
	id := "test-id"
	server := httptest.NewServer(newMockHandler(t, TestFileName, TestFileName, id))
	defer server.Close()

	pd := New()
	u, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal("Cannot parse a URL:", err)
	}
	pd.Client = client.NewHTTPClientWithConfig(nil, &client.TransportConfig{
		Host:     u.Host,
		BasePath: "/",
		Schemes:  []string{"http"},
	})

	fp, err := os.Open(TestFileName)
	if err != nil {
		t.Fatal("Failed to open the file:", err)
	}
	defer func() {
		if err := fp.Close(); err != nil && !errors.Is(err, os.ErrClosed) {
			t.Error(err)
		}
	}()

	res, err := pd.UploadRaw(context.Background(), fp, TestFileName)
	if err != nil {
		t.Fatal("failed to upload a file:", err.Error())
	}
	if res != id {
		t.Errorf("received ID = %v, want %v", res, id)
	}
}
