// Code generated by go-swagger; DO NOT EDIT.

package file

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"fmt"
	"io"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"

	"github.com/jkawamoto/go-pixeldrain/models"
)

// GetFileInfoReader is a Reader for the GetFileInfo structure.
type GetFileInfoReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetFileInfoReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetFileInfoOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewGetFileInfoDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetFileInfoOK creates a GetFileInfoOK with default headers values
func NewGetFileInfoOK() *GetFileInfoOK {
	return &GetFileInfoOK{}
}

/* GetFileInfoOK describes a response with status code 200, with default header values.

OK
*/
type GetFileInfoOK struct {
	Payload *GetFileInfoOKBody
}

func (o *GetFileInfoOK) Error() string {
	return fmt.Sprintf("[GET /file/{id}/info][%d] getFileInfoOK  %+v", 200, o.Payload)
}
func (o *GetFileInfoOK) GetPayload() *GetFileInfoOKBody {
	return o.Payload
}

func (o *GetFileInfoOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(GetFileInfoOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetFileInfoDefault creates a GetFileInfoDefault with default headers values
func NewGetFileInfoDefault(code int) *GetFileInfoDefault {
	return &GetFileInfoDefault{
		_statusCode: code,
	}
}

/* GetFileInfoDefault describes a response with status code -1, with default header values.

Error Response
*/
type GetFileInfoDefault struct {
	_statusCode int

	Payload *models.StandardError
}

// Code gets the status code for the get file info default response
func (o *GetFileInfoDefault) Code() int {
	return o._statusCode
}

func (o *GetFileInfoDefault) Error() string {
	return fmt.Sprintf("[GET /file/{id}/info][%d] getFileInfo default  %+v", o._statusCode, o.Payload)
}
func (o *GetFileInfoDefault) GetPayload() *models.StandardError {
	return o.Payload
}

func (o *GetFileInfoDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.StandardError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*GetFileInfoOKBody get file info o k body
swagger:model GetFileInfoOKBody
*/
type GetFileInfoOKBody struct {

	// Timestamp of last viewed time
	// Example: 2019-01-15T17:13:43Z
	// Format: date-time
	DateLastView strfmt.DateTime `json:"date_last_view,omitempty"`

	// Timestamp of uploaded time
	// Example: 2019-01-15T17:13:43Z
	// Format: date-time
	DateUpload strfmt.DateTime `json:"date_upload,omitempty"`

	// Description of the file
	// Example: File description
	Description string `json:"description,omitempty"`

	// ID of the newly uploaded file
	// Example: abc123
	ID string `json:"id,omitempty"`

	// Image associated with the mime type
	// Example: http://pixeldra.in/res/img/mime/image-png.png
	MimeImage string `json:"mime_image,omitempty"`

	// MIME type of the file
	// Example: image/png
	MimeType string `json:"mime_type,omitempty"`

	// Name of the file
	// Example: screenshot.png
	Name string `json:"name,omitempty"`

	// Size of the file in Bytes
	// Example: 5694837
	Size int64 `json:"size,omitempty"`

	// success
	// Example: true
	Success bool `json:"success,omitempty"`

	// Link to a thumbnail of this file
	// Example: http://pixeldra.in/api/thumbnail/123abc
	Thumbnail string `json:"thumbnail,omitempty"`

	// Amount of unique file views
	// Example: 1234
	Views int64 `json:"views,omitempty"`
}

// Validate validates this get file info o k body
func (o *GetFileInfoOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateDateLastView(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateDateUpload(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetFileInfoOKBody) validateDateLastView(formats strfmt.Registry) error {
	if swag.IsZero(o.DateLastView) { // not required
		return nil
	}

	if err := validate.FormatOf("getFileInfoOK"+"."+"date_last_view", "body", "date-time", o.DateLastView.String(), formats); err != nil {
		return err
	}

	return nil
}

func (o *GetFileInfoOKBody) validateDateUpload(formats strfmt.Registry) error {
	if swag.IsZero(o.DateUpload) { // not required
		return nil
	}

	if err := validate.FormatOf("getFileInfoOK"+"."+"date_upload", "body", "date-time", o.DateUpload.String(), formats); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this get file info o k body based on context it is used
func (o *GetFileInfoOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *GetFileInfoOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetFileInfoOKBody) UnmarshalBinary(b []byte) error {
	var res GetFileInfoOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}