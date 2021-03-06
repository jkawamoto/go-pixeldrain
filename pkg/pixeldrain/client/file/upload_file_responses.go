// Code generated by go-swagger; DO NOT EDIT.

package file

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/jkawamoto/go-pixeldrain/pkg/pixeldrain/models"
)

// UploadFileReader is a Reader for the UploadFile structure.
type UploadFileReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UploadFileReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 201:
		result := NewUploadFileCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewUploadFileDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewUploadFileCreated creates a UploadFileCreated with default headers values
func NewUploadFileCreated() *UploadFileCreated {
	return &UploadFileCreated{}
}

/* UploadFileCreated describes a response with status code 201, with default header values.

File is uploaded
*/
type UploadFileCreated struct {
	Payload *UploadFileCreatedBody
}

func (o *UploadFileCreated) Error() string {
	return fmt.Sprintf("[POST /file][%d] uploadFileCreated  %+v", 201, o.Payload)
}
func (o *UploadFileCreated) GetPayload() *UploadFileCreatedBody {
	return o.Payload
}

func (o *UploadFileCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(UploadFileCreatedBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUploadFileDefault creates a UploadFileDefault with default headers values
func NewUploadFileDefault(code int) *UploadFileDefault {
	return &UploadFileDefault{
		_statusCode: code,
	}
}

/* UploadFileDefault describes a response with status code -1, with default header values.

Error Response
*/
type UploadFileDefault struct {
	_statusCode int

	Payload *models.StandardError
}

// Code gets the status code for the upload file default response
func (o *UploadFileDefault) Code() int {
	return o._statusCode
}

func (o *UploadFileDefault) Error() string {
	return fmt.Sprintf("[POST /file][%d] uploadFile default  %+v", o._statusCode, o.Payload)
}
func (o *UploadFileDefault) GetPayload() *models.StandardError {
	return o.Payload
}

func (o *UploadFileDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.StandardError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*UploadFileCreatedBody upload file created body
swagger:model UploadFileCreatedBody
*/
type UploadFileCreatedBody struct {

	// ID of the newly uploaded file
	// Example: abc123
	ID string `json:"id,omitempty"`

	// success
	// Example: true
	Success bool `json:"success,omitempty"`
}

// Validate validates this upload file created body
func (o *UploadFileCreatedBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this upload file created body based on context it is used
func (o *UploadFileCreatedBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *UploadFileCreatedBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *UploadFileCreatedBody) UnmarshalBinary(b []byte) error {
	var res UploadFileCreatedBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
