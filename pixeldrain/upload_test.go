/*
 * upload_test.go
 *
 * Copyright (c) 2018-2019 Junpei Kawamoto
 *
 * This software is released under the MIT License.
 *
 * http://opensource.org/licenses/mit-license.php
 */

package pixeldrain

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"testing"

	"github.com/jkawamoto/go-pixeldrain/client"
	"github.com/jkawamoto/go-pixeldrain/client/file"
)

const (
	TestFileName = "upload.go"
)

type mockServer struct {
	expected []byte
	id       string
	name     string
}

func newMockServer(file, name, id string) (m *mockServer, err error) {

	fp, err := os.Open(file)
	if err != nil {
		return
	}
	//noinspection GoUnhandledErrorResult
	defer fp.Close()

	data, err := ioutil.ReadAll(fp)
	if err != nil {
		return
	}

	m = &mockServer{
		expected: data,
		id:       id,
		name:     name,
	}
	return

}

func (m *mockServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	if req.URL.Path != "/file" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if v := req.FormValue("name"); v != m.name {
		res.WriteHeader(http.StatusBadRequest)
		//noinspection GoUnhandledErrorResult
		fmt.Fprintf(res, "given file name is %v, want %v", v, m.name)
		return
	}

	fp, _, err := req.FormFile("file")
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		//noinspection GoUnhandledErrorResult
		fmt.Fprintln(res, "cannot get the file:", err.Error())
		return
	}
	//noinspection GoUnhandledErrorResult
	defer fp.Close()

	data, err := ioutil.ReadAll(fp)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		//noinspection GoUnhandledErrorResult
		fmt.Fprintln(res, "Cannot read the uploaded file:", err.Error())
		return
	}

	if !reflect.DeepEqual(data, m.expected) {
		res.WriteHeader(http.StatusInternalServerError)
		//noinspection GoUnhandledErrorResult
		fmt.Fprintln(res, "Uploaded file is broken")
		return
	}

	res.Header().Set(ContentType, "application/json")
	res.WriteHeader(http.StatusCreated)
	//noinspection GoUnhandledErrorResult
	json.NewEncoder(res).Encode(&file.UploadFileCreatedBody{
		ID:      m.id,
		Success: true,
	})

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
		{file: "../cmd/pd/command/upload.go", rename: "", expect: TestFileName},
	}

	for _, c := range cases {

		t.Run(fmt.Sprintf("%+v", c), func(t *testing.T) {

			m, err := newMockServer(c.file, c.expect, id)
			if err != nil {
				t.Fatal("Cannot prepare a mock server:", err)
			}
			server := httptest.NewServer(m)
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
			//noinspection GoUnhandledErrorResult
			defer fp.Close()

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
	m, err := newMockServer(TestFileName, TestFileName, id)
	if err != nil {
		t.Fatal("Cannot prepare a mock server:", err)
	}
	server := httptest.NewServer(m)
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
	//noinspection GoUnhandledErrorResult
	defer fp.Close()

	res, err := pd.UploadRaw(context.Background(), fp, TestFileName)
	if err != nil {
		t.Fatal("failed to upload a file:", err.Error())
	}
	if res != id {
		t.Errorf("received ID = %v, want %v", res, id)
	}

}
