// transporter.go
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
)

const ContentType = "content-type"

// transporter is a transporter which adds authentication information to each request before transporting it.
type transporter struct {
	http.RoundTripper
}

// newTransporter creates a transporter which wraps a given transporter.
func newTransporter(transport http.RoundTripper) (res *transporter) {

	res = &transporter{
		RoundTripper: transport,
	}
	return

}

func (t *transporter) RoundTrip(req *http.Request) (res *http.Response, err error) {

	res, err = t.RoundTripper.RoundTrip(req)
	if err != nil {
		return
	}

	if strings.HasPrefix(res.Header.Get(ContentType), "text/plain") {
		res.Header.Del(ContentType)
		res.Header.Set(ContentType, "application/json")
	}

	return

}
