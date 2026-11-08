// transport_test.go
//
// Copyright (c) 2018-2025 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package pixeldrain

import (
	"context"
	"net/http"
	"sync"
	"testing"

	"github.com/go-openapi/runtime"
)

type operationTransport struct{}

func (operationTransport) Submit(op *runtime.ClientOperation) (any, error) {
	return op, nil
}

func (operationTransport) SubmitContext(_ context.Context, op *runtime.ClientOperation) (any, error) {
	return op, nil
}

func TestSubmitDoesNotModifyOperation(t *testing.T) {
	base := &mockRoundTripper{Response: &http.Response{StatusCode: http.StatusOK, Header: make(http.Header)}}
	client := &http.Client{Transport: base}
	op := &runtime.ClientOperation{Client: client, ProducesMediaTypes: []string{runtime.JSONMime}}
	transport := ContentTypeFixer(operationTransport{})

	var wg sync.WaitGroup
	for range 10 {
		wg.Go(func() {
			for _, submit := range []func(*runtime.ClientOperation) (any, error){transport.Submit, func(op *runtime.ClientOperation) (any, error) {
				return transport.SubmitContext(context.Background(), op)
			}} {
				result, err := submit(op)
				if err != nil {
					t.Error(err)
					continue
				}
				prepared := result.(*runtime.ClientOperation)
				if prepared == op || prepared.Client == client {
					t.Error("operation and client must be copied")
				}
				if wrapper, ok := prepared.Client.Transport.(*roundTripper); !ok || wrapper.upstream != base || wrapper.contentType != runtime.JSONMime {
					t.Error("unexpected transport wrapper")
				}
			}
		})
	}
	wg.Wait()
	if op.Client != client || client.Transport != base {
		t.Error("original operation or client was modified")
	}
}

func TestSubmitWithoutMediaTypes(t *testing.T) {
	op := &runtime.ClientOperation{}
	result, err := ContentTypeFixer(operationTransport{}).Submit(op)
	if err != nil {
		t.Fatal(err)
	}
	prepared := result.(*runtime.ClientOperation)
	if prepared.Client == nil || prepared.Client.Transport == nil {
		t.Fatal("expected an initialized client and transport")
	}
	if _, ok := prepared.Client.Transport.(*roundTripper); ok {
		t.Error("content type must not be overridden without media types")
	}
	if op.Client != nil {
		t.Error("original operation was modified")
	}
}

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
