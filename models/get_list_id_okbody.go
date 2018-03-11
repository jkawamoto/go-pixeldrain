// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// GetListIDOKBody get list Id o k body
// swagger:model getListIdOKBody
type GetListIDOKBody struct {

	// date creqated
	DateCreqated float64 `json:"date_creqated,omitempty"`

	// files
	Files GetListIDOKBodyFiles `json:"files"`

	// id
	ID string `json:"id,omitempty"`

	// success
	Success bool `json:"success,omitempty"`

	// title
	Title string `json:"title,omitempty"`
}

// Validate validates this get list Id o k body
func (m *GetListIDOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *GetListIDOKBody) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *GetListIDOKBody) UnmarshalBinary(b []byte) error {
	var res GetListIDOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
