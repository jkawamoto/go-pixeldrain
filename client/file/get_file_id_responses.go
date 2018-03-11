// Code generated by go-swagger; DO NOT EDIT.

package file

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/jkawamoto/go-pixeldrain/models"
)

// GetFileIDReader is a Reader for the GetFileID structure.
type GetFileIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetFileIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewGetFileIDOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		result := NewGetFileIDDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetFileIDOK creates a GetFileIDOK with default headers values
func NewGetFileIDOK() *GetFileIDOK {
	return &GetFileIDOK{}
}

/*GetFileIDOK handles this case with default header values.

A file output stream.
*/
type GetFileIDOK struct {
	Payload models.GetFileIDOKBody
}

func (o *GetFileIDOK) Error() string {
	return fmt.Sprintf("[GET /file/{id}][%d] getFileIdOK  %+v", 200, o.Payload)
}

func (o *GetFileIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetFileIDDefault creates a GetFileIDDefault with default headers values
func NewGetFileIDDefault(code int) *GetFileIDDefault {
	return &GetFileIDDefault{
		_statusCode: code,
	}
}

/*GetFileIDDefault handles this case with default header values.

Error Response
*/
type GetFileIDDefault struct {
	_statusCode int

	Payload *models.StandardError
}

// Code gets the status code for the get file ID default response
func (o *GetFileIDDefault) Code() int {
	return o._statusCode
}

func (o *GetFileIDDefault) Error() string {
	return fmt.Sprintf("[GET /file/{id}][%d] GetFileID default  %+v", o._statusCode, o.Payload)
}

func (o *GetFileIDDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.StandardError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
