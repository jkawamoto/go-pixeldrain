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

	"github.com/jkawamoto/go-pixeldrain/pkg/pixeldrain/models"
)

type Error struct {
	error
}

func NewError(err error) *Error {
	if err == nil {
		return nil
	}
	return &Error{error: err}
}

func (e Error) Error() string {
	if p, ok := e.error.(interface {
		GetPayload() *models.StandardError
	}); ok {
		return swag.StringValue(p.GetPayload().Message)
	}
	return e.error.Error()
}

func (e Error) Unwrap() error {
	return e.error
}
