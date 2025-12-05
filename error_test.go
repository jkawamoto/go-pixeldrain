// error_test.go
//
// Copyright (c) 2018-2025 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package pixeldrain

import (
	"context"
	"errors"
	"testing"

	"github.com/go-openapi/swag"

	"github.com/jkawamoto/go-pixeldrain/models"
)

type apiError struct {
	err *models.StandardError
}

func newAPIError(e *models.StandardError) error {
	return &apiError{err: e}
}

func (e *apiError) Error() string {
	return "unexpected call"
}

func (e *apiError) GetPayload() *models.StandardError {
	return e.err
}

func TestNewError(t *testing.T) {
	sampleMsg := "this is a sample error message"

	cases := []struct {
		name   string
		err    error
		expect string
	}{
		{
			name: "API error",
			err: newAPIError(&models.StandardError{
				Message: swag.String(sampleMsg),
			}),
			expect: sampleMsg,
		},
		{
			name:   "non API error",
			err:    context.Canceled,
			expect: context.Canceled.Error(),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := NewError(c.err)
			if msg := err.Error(); msg != c.expect {
				t.Errorf("expect %v, got %v", c.expect, msg)
			}
			if !errors.Is(err, c.err) {
				t.Errorf("expect %v is a %v", c.err, err)
			}
		})
	}
}
