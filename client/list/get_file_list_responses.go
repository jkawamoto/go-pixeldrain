// Code generated by go-swagger; DO NOT EDIT.

package list

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/jkawamoto/go-pixeldrain/models"
)

// GetFileListReader is a Reader for the GetFileList structure.
type GetFileListReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetFileListReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewGetFileListOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		result := NewGetFileListDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetFileListOK creates a GetFileListOK with default headers values
func NewGetFileListOK() *GetFileListOK {
	return &GetFileListOK{}
}

/*GetFileListOK handles this case with default header values.

OK
*/
type GetFileListOK struct {
	Payload *GetFileListOKBody
}

func (o *GetFileListOK) Error() string {
	return fmt.Sprintf("[GET /list/{id}][%d] getFileListOK  %+v", 200, o.Payload)
}

func (o *GetFileListOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(GetFileListOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetFileListDefault creates a GetFileListDefault with default headers values
func NewGetFileListDefault(code int) *GetFileListDefault {
	return &GetFileListDefault{
		_statusCode: code,
	}
}

/*GetFileListDefault handles this case with default header values.

Error Response
*/
type GetFileListDefault struct {
	_statusCode int

	Payload *models.StandardError
}

// Code gets the status code for the get file list default response
func (o *GetFileListDefault) Code() int {
	return o._statusCode
}

func (o *GetFileListDefault) Error() string {
	return fmt.Sprintf("[GET /list/{id}][%d] getFileList default  %+v", o._statusCode, o.Payload)
}

func (o *GetFileListDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.StandardError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*GetFileListOKBody get file list o k body
swagger:model GetFileListOKBody
*/
type GetFileListOKBody struct {

	// date creqated
	DateCreqated float64 `json:"date_creqated,omitempty"`

	// files
	Files []*FilesItems0 `json:"files"`

	// id
	ID string `json:"id,omitempty"`

	// success
	Success bool `json:"success,omitempty"`

	// title
	Title string `json:"title,omitempty"`
}

// Validate validates this get file list o k body
func (o *GetFileListOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateFiles(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetFileListOKBody) validateFiles(formats strfmt.Registry) error {

	if swag.IsZero(o.Files) { // not required
		return nil
	}

	for i := 0; i < len(o.Files); i++ {
		if swag.IsZero(o.Files[i]) { // not required
			continue
		}

		if o.Files[i] != nil {
			if err := o.Files[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getFileListOK" + "." + "files" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetFileListOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetFileListOKBody) UnmarshalBinary(b []byte) error {
	var res GetFileListOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
