// pixeldrain.go
//
// Copyright (c) 2018-2025 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package pixeldrain

import (
	"net/url"
	"strings"

	"github.com/go-openapi/strfmt"

	"github.com/jkawamoto/go-pixeldrain/client"
)

// Default is a default client.
var Default = New(nil, nil)

const (
	// fileBasePath is the base path to file download URLs.
	fileBasePath = "/file"
	// listBasePath is the base path to list URLs.
	listBasePath = "/l"
)

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
	return u.JoinPath(client.DefaultHost, client.DefaultBasePath, fileBasePath, id).String()
}

// ListURL returns the URL associated with the given list ID.
func ListURL(id string) string {
	u := url.URL{
		Scheme: client.DefaultSchemes[0],
	}
	return u.JoinPath(client.DefaultHost, listBasePath, id).String()
}

// IsDownloadURL returns true if the given url points a file.
func IsDownloadURL(u string) (bool, error) {
	parse, err := url.Parse(u)
	if err != nil {
		return false, err
	}

	prefix, err := url.JoinPath(client.DefaultBasePath, fileBasePath)
	if err != nil {
		return false, err
	}

	return strings.HasPrefix(parse.Path, prefix), nil
}

// IsListURL returns true if the given url points a list.
func IsListURL(u string) (bool, error) {
	parse, err := url.Parse(u)
	if err != nil {
		return false, err
	}

	return strings.HasPrefix(parse.Path, listBasePath), nil
}
