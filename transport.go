// transport.go
//
// Copyright (c) 2018-2025 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package pixeldrain

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// roundTripper is a http.RoundTripper that forwards a request to the upstream and fixes content type header of the
// corresponding response.
type roundTripper struct {
	upstream    http.RoundTripper
	contentType string
}

var _ http.RoundTripper = (*roundTripper)(nil)

// newRoundTripper creates a roundTripper which wraps a given roundTripper and overwrites the content types of responses.
func newRoundTripper(upstream http.RoundTripper, contentType string) *roundTripper {
	return &roundTripper{
		upstream:    upstream,
		contentType: contentType,
	}
}

// RoundTrip executes a single HTTP transaction, returning a Response for the provided Request.
func (t *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	res, err := t.upstream.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 300 {
		res.Header.Set(runtime.HeaderContentType, t.contentType)
	}
	return res, nil
}

// transport is a runtime.ClientTransport that modifies the http client of each request and forwards the request to the
// upstream transport.
type transport struct {
	upstream runtime.ClientTransport
}

var _ runtime.ClientTransport = (*transport)(nil)

// ContentTypeFixer returns a new ClientTransport that wraps the given ClientTransport to fix the content type issue.
func ContentTypeFixer(upstream runtime.ClientTransport) runtime.ClientTransport {
	return &transport{upstream: upstream}
}

// Submit sends the given operation and returns a response.
func (t *transport) Submit(op *runtime.ClientOperation) (interface{}, error) {
	if op.Client == nil {
		op.Client = &http.Client{
			Transport: http.DefaultTransport,
		}
	}
	op.Client.Transport = newRoundTripper(op.Client.Transport, op.ProducesMediaTypes[0])
	return t.upstream.Submit(op)
}
