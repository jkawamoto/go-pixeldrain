// The MIT License (MIT)
//
// Copyright (c) 2018 Junpei Kawamoto
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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
		fmt.Fprintf(res, "given file name is %v, want %v", v, m.name)
		return
	}

	fp, _, err := req.FormFile("file")
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(res, "cannot get the file:", err.Error())
		return
	}
	defer fp.Close()

	data, err := ioutil.ReadAll(fp)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(res, "Cannot read the uploaded file:", err.Error())
		return
	}

	if !reflect.DeepEqual(data, m.expected) {
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(res, "Uploaded file is broken")
		return
	}

	res.WriteHeader(http.StatusCreated)
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
		{file: "../command/upload.go", rename: "", expect: TestFileName},
	}

	for _, c := range cases {

		t.Run(fmt.Sprintf("%+v", c), func(t *testing.T) {

			m, err := newMockServer(c.file, c.expect, id)
			if err != nil {
				t.Fatal("Cannot prepare a mock server:", err)
			}
			server := httptest.NewServer(m)
			defer server.Close()

			u, err := url.Parse(server.URL)
			if err != nil {
				t.Fatal("Cannot parse a URL:", err)
			}
			cli := client.NewHTTPClientWithConfig(nil, &client.TransportConfig{
				Host:     u.Host,
				BasePath: "/",
				Schemes:  []string{"http"},
			})

			fp, err := os.Open(c.file)
			if err != nil {
				t.Fatal("Failed to open the file:", err)
			}
			defer fp.Close()

			res, err := Upload(context.Background(), cli, fp, c.rename)
			if err != nil {
				t.Fatal("failed to upload a file:", err.Error())
			}
			if res != id {
				t.Errorf("received ID = %v, want %v", res, id)
			}

		})

	}

}
