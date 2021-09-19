// transport.go
//
// Copyright (c) 2018-2021 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package pixeldrain

import (
	"net/http"
	"strings"

	"github.com/go-openapi/runtime"
)

const ContentType = "content-type"

// roundTripper is a http.RoundTripper that forwards a request to the upstream and fixes content type header of the
// corresponding response.
type roundTripper struct {
	upstream http.RoundTripper
}

var _ http.RoundTripper = (*roundTripper)(nil)

// newRoundTripper creates a roundTripper which wraps a given roundTripper.
func newRoundTripper(upstream http.RoundTripper) *roundTripper {
	return &roundTripper{
		upstream: upstream,
	}
}

// RoundTrip executes a single HTTP transaction, returning a Response for the provided Request.
func (t *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	res, err := t.upstream.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(res.Header.Get(ContentType), "text/plain") {
		res.Header.Set(ContentType, "application/json")
	}

	return res, nil
}

// transport is a runtime.ClientTransport that modifies the http client of each request and forwards the request to the
// upstream transport.
type transport struct {
	upstream runtime.ClientTransport
}

var _ runtime.ClientTransport = (*transport)(nil)

// newTransport returns a new ClientTransport that wraps the given ClientTransport.
func newTransport(upstream runtime.ClientTransport) runtime.ClientTransport {
	return &transport{upstream: upstream}
}

func (t *transport) Submit(op *runtime.ClientOperation) (interface{}, error) {
	if op.Client == nil {
		op.Client = &http.Client{
			Transport: http.DefaultTransport,
		}
	}
	op.Client.Transport = newRoundTripper(op.Client.Transport)
	return t.upstream.Submit(op)
}
