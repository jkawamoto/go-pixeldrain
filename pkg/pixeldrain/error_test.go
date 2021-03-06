// error_test.go
//
// Copyright (c) 2018-2021 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package pixeldrain

import (
	"testing"

	"github.com/go-openapi/swag"

	"github.com/jkawamoto/go-pixeldrain/pkg/pixeldrain/models"
)

type errorResponse struct {
	Err *models.StandardError
}

func (e *errorResponse) GetPayload() *models.StandardError {
	return e.Err
}

func TestNewAPIError(t *testing.T) {
	cases := []struct {
		name   string
		res    *errorResponse
		expect string
	}{
		{
			name: "error with a message",
			res: &errorResponse{
				Err: &models.StandardError{
					Message: swag.String("a message"),
				},
			},
			expect: "a message",
		},
		{
			name: "error without a message",
			res: &errorResponse{
				Err: &models.StandardError{},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := NewAPIError(c.res)
			if err == nil {
				t.Fatal("expect a non nil error")
			}

			expect := "API error: " + c.expect
			if err.Error() != expect {
				t.Errorf("expect %q, got %q", expect, err.Error())
			}
		})
	}
}
