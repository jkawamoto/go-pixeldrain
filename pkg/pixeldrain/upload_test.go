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

	"github.com/go-openapi/swag"

	"github.com/jkawamoto/go-pixeldrain/pkg/pixeldrain/client"
	"github.com/jkawamoto/go-pixeldrain/pkg/pixeldrain/client/file"
	"github.com/jkawamoto/go-pixeldrain/pkg/pixeldrain/models"
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
	res.Header().Set(ContentType, "application/json")
	if req.URL.Path != "/file" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if v := req.FormValue("name"); v != m.name {
		res.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(res).Encode(&models.StandardError{
			Message: swag.String(fmt.Sprintf("given file name is %v, want %v", v, m.name)),
		}); err != nil {
			panic(err)
		}
		return
	}

	fp, _, err := req.FormFile("file")
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(res).Encode(&models.StandardError{
			Message: swag.String(fmt.Sprint("cannot get the file:", err)),
		}); err != nil {
			panic(err)
		}
		return
	}
	defer func() {
		if err := fp.Close(); err != nil && !errors.Is(err, os.ErrClosed) {
			panic(err)
		}
	}()

	data, err := ioutil.ReadAll(fp)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(res).Encode(&models.StandardError{
			Message: swag.String(fmt.Sprint("Cannot read the uploaded file:", err)),
		}); err != nil {
			panic(err)
		}
		return
	}

	if !reflect.DeepEqual(data, m.expected) {
		res.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(res).Encode(&models.StandardError{
			Message: swag.String("Uploaded file is broken"),
		}); err != nil {
			panic(err)
		}
		return
	}

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
		name   string
		expect string
	}{
		{file: TestFileName, name: "", expect: TestFileName},
		{file: TestFileName, name: "another-expect", expect: "another-expect"},
		{file: "../../cmd/pd/command/upload.go", name: "", expect: TestFileName},
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

			res, err := pd.Upload(context.Background(), fp, c.name)
			if err != nil {
				t.Fatal("failed to upload a file:", err)
			}
			if res != id {
				t.Errorf("received ID = %v, want %v", res, id)
			}
		})
	}
}
