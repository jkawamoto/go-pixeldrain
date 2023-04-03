// pixeldrain.go
//
// Copyright (c) 2018-2023 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package pixeldrain

import (
	"net/url"

	"github.com/go-openapi/strfmt"

	"github.com/jkawamoto/go-pixeldrain/client"
)

// Default is a default client.
var Default = New(nil, nil)

// New creates a new PixelDrain client with the given configurations. formats and cfg can be nil.
func New(formats strfmt.Registry, cfg *client.TransportConfig) *client.PixeldrainAPI {
	cli := client.NewHTTPClientWithConfig(formats, cfg)
	cli.SetTransport(ContentTypeFixer(cli.Transport))
	return cli
}

// DownloadURL returns the URL associated with the given file ID.
func DownloadURL(id string) string {
	u := url.URL{
		Scheme: client.DefaultSchemes[0],
	}
	return u.JoinPath(client.DefaultHost, client.DefaultBasePath, "file", id).String()
}

// ListURL returns the URL associated with the given list ID.
func ListURL(id string) string {
	u := url.URL{
		Scheme: client.DefaultSchemes[0],
	}
	return u.JoinPath(client.DefaultHost, "l", id).String()
}
