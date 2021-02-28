// transporter_test.go
//
// Copyright (c) 2018-2021 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package pixeldrain

import (
	"net/http"
	"testing"
)

type mockRoundTripper struct {
	Request  *http.Request
	Response *http.Response
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	m.Request = req
	return m.Response, nil
}

func TestTransporter(t *testing.T) {
	cases := []struct {
		receive string
		expect  string
	}{
		{receive: "text/plain charset=utf8", expect: "application/json"},
		{receive: "image/png", expect: "image/png"},
	}

	for _, c := range cases {

		t.Run(c.receive, func(t *testing.T) {

			mock := &mockRoundTripper{
				Response: &http.Response{
					Status:     "OK",
					StatusCode: http.StatusOK,
					Header:     make(http.Header),
				},
			}
			mock.Response.Header.Set(ContentType, c.receive)

			transporter := newTransporter(mock)
			res, err := transporter.RoundTrip(nil)
			if err != nil {
				t.Fatal("RoundTrip returns an error:", err)
			}
			if cType := res.Header.Get(ContentType); cType != c.expect {
				t.Errorf("content-type is %v, expected %v", cType, c.expect)
			}

		})
	}
}
