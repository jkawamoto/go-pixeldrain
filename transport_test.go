// transport_test.go
//
// Copyright (c) 2018-2025 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package pixeldrain

import (
	"net/http"
	"testing"

	"github.com/go-openapi/runtime"
)

type mockRoundTripper struct {
	Request  *http.Request
	Response *http.Response
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	m.Request = req
	return m.Response, nil
}

func TestContentTypeFixer(t *testing.T) {
	cases := []struct {
		name        string
		status      int
		contentType string
		expect      string
	}{
		{
			name:        "replace with JSON",
			status:      http.StatusOK,
			contentType: runtime.JSONMime,
			expect:      runtime.JSONMime,
		},
		{
			name:        "replace with octet-stream",
			status:      http.StatusOK,
			contentType: runtime.DefaultMime,
			expect:      runtime.DefaultMime,
		},
		{
			name:        "error",
			status:      http.StatusNotFound,
			contentType: runtime.DefaultMime,
			expect:      "text/plain charset=utf8",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mock := &mockRoundTripper{
				Response: &http.Response{
					Status:     "OK",
					StatusCode: c.status,
					Header:     make(http.Header),
				},
			}
			mock.Response.Header.Set(runtime.HeaderContentType, "text/plain charset=utf8")

			transporter := newRoundTripper(mock, c.contentType)
			res, err := transporter.RoundTrip(nil)
			if err != nil {
				t.Fatal("RoundTrip returns an error:", err)
			}
			if cType := res.Header.Get(runtime.HeaderContentType); cType != c.expect {
				t.Errorf("content-type is %v, expected %v", cType, c.expect)
			}
		})
	}
}

func TestOptimizeTransport(t *testing.T) {
	t.Run("nil transport returns optimized transport", func(t *testing.T) {
		rt := optimizeTransport(nil)
		tr, ok := rt.(*http.Transport)
		if !ok {
			t.Fatalf("expected *http.Transport, got %T", rt)
		}
		if tr.WriteBufferSize != BufferSize {
			t.Errorf("expected WriteBufferSize %d, got %d", BufferSize, tr.WriteBufferSize)
		}
		if tr.ReadBufferSize != BufferSize {
			t.Errorf("expected ReadBufferSize %d, got %d", BufferSize, tr.ReadBufferSize)
		}
	})

	t.Run("existing http.Transport gets enlarged buffers", func(t *testing.T) {
		base := &http.Transport{}
		rt := optimizeTransport(base)
		tr, ok := rt.(*http.Transport)
		if !ok {
			t.Fatalf("expected *http.Transport, got %T", rt)
		}
		if tr.WriteBufferSize != BufferSize {
			t.Errorf("expected WriteBufferSize %d, got %d", BufferSize, tr.WriteBufferSize)
		}
		if tr.ReadBufferSize != BufferSize {
			t.Errorf("expected ReadBufferSize %d, got %d", BufferSize, tr.ReadBufferSize)
		}
		// ensure original was not mutated
		if base.WriteBufferSize != 0 {
			t.Errorf("expected original WriteBufferSize 0, got %d", base.WriteBufferSize)
		}
	})

	t.Run("custom RoundTripper passed through", func(t *testing.T) {
		mock := &mockRoundTripper{}
		rt := optimizeTransport(mock)
		if rt != mock {
			t.Errorf("expected mock to be passed through, got %T", rt)
		}
	})
}
