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

	"github.com/jkawamoto/go-pixeldrain/models"
)

// Error wraps an error to return a more readable error message.
type Error struct {
	error
}

// NewError returns *Error that wraps the given error. If err is nil, it returns nil.
func NewError(err error) *Error {
	if err == nil {
		return nil
	}
	return &Error{error: err}
}

// Error returns a string that represents this error.
func (e Error) Error() string {
	if p, ok := e.error.(interface {
		GetPayload() *models.StandardError
	}); ok {
		return swag.StringValue(p.GetPayload().Message)
	}
	return e.error.Error()
}

// Unwrap returns the original error.
func (e Error) Unwrap() error {
	return e.error
}
