// download_test.go
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
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-openapi/runtime"

	"github.com/jkawamoto/go-pixeldrain/pkg/pixeldrain/client"
	"github.com/jkawamoto/go-pixeldrain/pkg/pixeldrain/client/file"
)

type mockDownloadHandler struct {
	ID            string
	File          string
	Authorization string
}

func (m *mockDownloadHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(runtime.HeaderContentType, runtime.JSONMime)
	if !strings.HasPrefix(req.URL.Path, "/file") {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.Header.Get("Authorization") != m.Authorization {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !strings.Contains(req.URL.Path, m.ID) {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	if strings.HasSuffix(req.URL.Path, "/info") {
		info, err := os.Stat(m.File)
		if err != nil {
			res.WriteHeader(http.StatusNotFound)
			return
		}
		res.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(res).Encode(&file.GetFileInfoOKBody{
			Success: true,
			ID:      m.ID,
			Name:    m.File,
			Size:    info.Size(),
		}); err != nil {
			panic(err)
		}
	} else {
		fp, err := os.Open(m.File)
		if err != nil {
			res.WriteHeader(http.StatusNotFound)
			return
		}
		defer func() {
			if err := fp.Close(); err != nil && !errors.Is(err, os.ErrClosed) {
				panic(err)
			}
		}()

		res.Header().Set(runtime.HeaderContentType, runtime.DefaultMime)
		res.WriteHeader(http.StatusOK)
		if _, err := io.Copy(res, fp); err != nil {
			panic(err)
		}
	}
}

func TestDownload(t *testing.T) {
	id := "abcde"
	filename := "./download.go"
	apiKey := "test api key"

	t.Run("stdout", func(t *testing.T) {
		server := httptest.NewServer(&mockDownloadHandler{
			ID:            id,
			File:          filename,
			Authorization: authorization(apiKey),
		})
		defer server.Close()

		pd := New(apiKey)
		u, err := url.Parse(server.URL)
		if err != nil {
			t.Fatal("Cannot parse a URL:", err)
		}
		pd.cli = client.NewHTTPClientWithConfig(nil, &client.TransportConfig{
			Host:     u.Host,
			BasePath: "/",
			Schemes:  []string{"http"},
		})

		tmp, err := ioutil.TempFile("", "")
		if err != nil {
			t.Fatal("Failed to create a temporal filename", err)
		}
		pd.Stdout = tmp
		defer func() {
			if err := tmp.Close(); err != nil {
				t.Error(err)
			}
		}()

		err = pd.Download(context.Background(), pd.DownloadURL(id), "")
		if err != nil {
			t.Fatal("failed to download the filename:", err)
		}

		received, err := ioutil.ReadFile(tmp.Name())
		if err != nil {
			t.Fatal("Failed to read received filename:", err)
		}
		expected, err := ioutil.ReadFile(filename)
		if err != nil {
			t.Fatal("Failed to read original filename:", err)
		}
		if string(received) != string(expected) {
			t.Error("Downloaded filename is broken")
		}
	})

	t.Run("dir", func(t *testing.T) {
		server := httptest.NewServer(&mockDownloadHandler{
			ID:            id,
			File:          filename,
			Authorization: authorization(apiKey),
		})
		defer server.Close()

		pd := New(apiKey)
		u, err := url.Parse(server.URL)
		if err != nil {
			t.Fatal("Cannot parse a URL:", err)
		}
		pd.cli = client.NewHTTPClientWithConfig(nil, &client.TransportConfig{
			Host:     u.Host,
			BasePath: "/",
			Schemes:  []string{"http"},
		})

		tmp, err := ioutil.TempDir("", "")
		if err != nil {
			t.Fatal("Failed to create a temporal directory", err)
		}

		err = pd.Download(context.Background(), pd.DownloadURL(id), tmp)
		if err != nil {
			t.Fatal("failed to download the filename:", err)
		}

		received, err := ioutil.ReadFile(filepath.Join(tmp, filename))
		if err != nil {
			t.Fatal("Failed to read received filename:", err)
		}
		expected, err := ioutil.ReadFile(filename)
		if err != nil {
			t.Fatal("Failed to read original filename:", err)
		}
		if string(received) != string(expected) {
			t.Error("Downloaded file is broken")
		}
	})
}
