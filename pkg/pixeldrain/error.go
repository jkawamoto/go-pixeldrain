// error.go
//
// Copyright (c) 2018-2021 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package pixeldrain

import (
	"github.com/go-openapi/swag"
	"github.com/pkg/errors"

	"github.com/jkawamoto/go-pixeldrain/pkg/pixeldrain/models"
)

// ErrorResponse is an interface a Swagger's error response has
type ErrorResponse interface {
	GetPayload() *models.StandardError
}

// NewAPIError converts the error response to an error.
func NewAPIError(e ErrorResponse) error {
	return errors.New("API error: " + swag.StringValue(e.GetPayload().Message))
}
