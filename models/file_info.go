// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// FileInfo file info
//
// swagger:model FileInfo
type FileInfo struct {

	// How much bandwidth this file used
	BandwidthUsed int64 `json:"bandwidth_used,omitempty"`

	// Timestamp of last viewed time
	// Example: 2019-01-15T17:13:43Z
	// Format: date-time
	DateLastView strfmt.DateTime `json:"date_last_view,omitempty"`

	// Timestamp of uploaded time
	// Example: 2019-01-15T17:13:43Z
	// Format: date-time
	DateUpload strfmt.DateTime `json:"date_upload,omitempty"`

	// ID of the newly uploaded file
	// Example: abc123
	// Required: true
	ID *string `json:"id"`

	// MIME type of the file
	// Example: image/png
	MimeType string `json:"mime_type,omitempty"`

	// Name of the file
	// Example: screenshot.png
	Name string `json:"name,omitempty"`

	// Size of the file in Bytes
	// Example: 5694837
	Size int64 `json:"size,omitempty"`

	// Link to a thumbnail of this file
	// Example: /file/1234abcd/thumbnail
	ThumbnailHref string `json:"thumbnail_href,omitempty"`

	// Amount of unique file views
	// Example: 1234
	Views int64 `json:"views,omitempty"`
}

// Validate validates this file info
func (m *FileInfo) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDateLastView(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDateUpload(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *FileInfo) validateDateLastView(formats strfmt.Registry) error {
	if swag.IsZero(m.DateLastView) { // not required
		return nil
	}

	if err := validate.FormatOf("date_last_view", "body", "date-time", m.DateLastView.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *FileInfo) validateDateUpload(formats strfmt.Registry) error {
	if swag.IsZero(m.DateUpload) { // not required
		return nil
	}

	if err := validate.FormatOf("date_upload", "body", "date-time", m.DateUpload.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *FileInfo) validateID(formats strfmt.Registry) error {

	if err := validate.Required("id", "body", m.ID); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this file info based on context it is used
func (m *FileInfo) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *FileInfo) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *FileInfo) UnmarshalBinary(b []byte) error {
	var res FileInfo
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}