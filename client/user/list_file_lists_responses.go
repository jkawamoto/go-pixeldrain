// Code generated by go-swagger; DO NOT EDIT.

package user

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/jkawamoto/go-pixeldrain/models"
)

// ListFileListsReader is a Reader for the ListFileLists structure.
type ListFileListsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ListFileListsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewListFileListsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewListFileListsDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewListFileListsOK creates a ListFileListsOK with default headers values
func NewListFileListsOK() *ListFileListsOK {
	return &ListFileListsOK{}
}

/*
ListFileListsOK describes a response with status code 200, with default header values.

OK
*/
type ListFileListsOK struct {
	Payload *ListFileListsOKBody
}

// IsSuccess returns true when this list file lists o k response has a 2xx status code
func (o *ListFileListsOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this list file lists o k response has a 3xx status code
func (o *ListFileListsOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this list file lists o k response has a 4xx status code
func (o *ListFileListsOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this list file lists o k response has a 5xx status code
func (o *ListFileListsOK) IsServerError() bool {
	return false
}

// IsCode returns true when this list file lists o k response a status code equal to that given
func (o *ListFileListsOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the list file lists o k response
func (o *ListFileListsOK) Code() int {
	return 200
}

func (o *ListFileListsOK) Error() string {
	return fmt.Sprintf("[GET /user/lists][%d] listFileListsOK  %+v", 200, o.Payload)
}

func (o *ListFileListsOK) String() string {
	return fmt.Sprintf("[GET /user/lists][%d] listFileListsOK  %+v", 200, o.Payload)
}

func (o *ListFileListsOK) GetPayload() *ListFileListsOKBody {
	return o.Payload
}

func (o *ListFileListsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(ListFileListsOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewListFileListsDefault creates a ListFileListsDefault with default headers values
func NewListFileListsDefault(code int) *ListFileListsDefault {
	return &ListFileListsDefault{
		_statusCode: code,
	}
}

/*
ListFileListsDefault describes a response with status code -1, with default header values.

Error Response
*/
type ListFileListsDefault struct {
	_statusCode int

	Payload *models.StandardError
}

// IsSuccess returns true when this list file lists default response has a 2xx status code
func (o *ListFileListsDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this list file lists default response has a 3xx status code
func (o *ListFileListsDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this list file lists default response has a 4xx status code
func (o *ListFileListsDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this list file lists default response has a 5xx status code
func (o *ListFileListsDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this list file lists default response a status code equal to that given
func (o *ListFileListsDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the list file lists default response
func (o *ListFileListsDefault) Code() int {
	return o._statusCode
}

func (o *ListFileListsDefault) Error() string {
	return fmt.Sprintf("[GET /user/lists][%d] listFileLists default  %+v", o._statusCode, o.Payload)
}

func (o *ListFileListsDefault) String() string {
	return fmt.Sprintf("[GET /user/lists][%d] listFileLists default  %+v", o._statusCode, o.Payload)
}

func (o *ListFileListsDefault) GetPayload() *models.StandardError {
	return o.Payload
}

func (o *ListFileListsDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.StandardError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*
ListFileListsOKBody list file lists o k body
swagger:model ListFileListsOKBody
*/
type ListFileListsOKBody struct {

	// lists
	Lists []*models.ListInfo `json:"lists"`
}

// Validate validates this list file lists o k body
func (o *ListFileListsOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateLists(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *ListFileListsOKBody) validateLists(formats strfmt.Registry) error {
	if swag.IsZero(o.Lists) { // not required
		return nil
	}

	for i := 0; i < len(o.Lists); i++ {
		if swag.IsZero(o.Lists[i]) { // not required
			continue
		}

		if o.Lists[i] != nil {
			if err := o.Lists[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("listFileListsOK" + "." + "lists" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("listFileListsOK" + "." + "lists" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this list file lists o k body based on the context it is used
func (o *ListFileListsOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateLists(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *ListFileListsOKBody) contextValidateLists(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(o.Lists); i++ {

		if o.Lists[i] != nil {
			if err := o.Lists[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("listFileListsOK" + "." + "lists" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("listFileListsOK" + "." + "lists" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *ListFileListsOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *ListFileListsOKBody) UnmarshalBinary(b []byte) error {
	var res ListFileListsOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
