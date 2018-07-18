// Code generated by go-swagger; DO NOT EDIT.

package file

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"
)

// NewUploadFileParams creates a new UploadFileParams object
// with the default values initialized.
func NewUploadFileParams() *UploadFileParams {
	var (
		descriptionDefault = string("Pixeldrain File")
		nameDefault        = string("Name of file param")
	)
	return &UploadFileParams{
		Description: &descriptionDefault,
		Name:        &nameDefault,

		timeout: cr.DefaultTimeout,
	}
}

// NewUploadFileParamsWithTimeout creates a new UploadFileParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewUploadFileParamsWithTimeout(timeout time.Duration) *UploadFileParams {
	var (
		descriptionDefault = string("Pixeldrain File")
		nameDefault        = string("Name of file param")
	)
	return &UploadFileParams{
		Description: &descriptionDefault,
		Name:        &nameDefault,

		timeout: timeout,
	}
}

// NewUploadFileParamsWithContext creates a new UploadFileParams object
// with the default values initialized, and the ability to set a context for a request
func NewUploadFileParamsWithContext(ctx context.Context) *UploadFileParams {
	var (
		descriptionDefault = string("Pixeldrain File")
		nameDefault        = string("Name of file param")
	)
	return &UploadFileParams{
		Description: &descriptionDefault,
		Name:        &nameDefault,

		Context: ctx,
	}
}

// NewUploadFileParamsWithHTTPClient creates a new UploadFileParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewUploadFileParamsWithHTTPClient(client *http.Client) *UploadFileParams {
	var (
		descriptionDefault = string("Pixeldrain File")
		nameDefault        = string("Name of file param")
	)
	return &UploadFileParams{
		Description: &descriptionDefault,
		Name:        &nameDefault,
		HTTPClient:  client,
	}
}

/*UploadFileParams contains all the parameters to send to the API endpoint
for the upload file operation typically these are written to a http.Request
*/
type UploadFileParams struct {

	/*Description
	  Description of the file

	*/
	Description *string
	/*File
	  Multipart file to upload

	*/
	File runtime.NamedReadCloser
	/*Name
	  Name of the file to upload

	*/
	Name *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the upload file params
func (o *UploadFileParams) WithTimeout(timeout time.Duration) *UploadFileParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the upload file params
func (o *UploadFileParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the upload file params
func (o *UploadFileParams) WithContext(ctx context.Context) *UploadFileParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the upload file params
func (o *UploadFileParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the upload file params
func (o *UploadFileParams) WithHTTPClient(client *http.Client) *UploadFileParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the upload file params
func (o *UploadFileParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithDescription adds the description to the upload file params
func (o *UploadFileParams) WithDescription(description *string) *UploadFileParams {
	o.SetDescription(description)
	return o
}

// SetDescription adds the description to the upload file params
func (o *UploadFileParams) SetDescription(description *string) {
	o.Description = description
}

// WithFile adds the file to the upload file params
func (o *UploadFileParams) WithFile(file runtime.NamedReadCloser) *UploadFileParams {
	o.SetFile(file)
	return o
}

// SetFile adds the file to the upload file params
func (o *UploadFileParams) SetFile(file runtime.NamedReadCloser) {
	o.File = file
}

// WithName adds the name to the upload file params
func (o *UploadFileParams) WithName(name *string) *UploadFileParams {
	o.SetName(name)
	return o
}

// SetName adds the name to the upload file params
func (o *UploadFileParams) SetName(name *string) {
	o.Name = name
}

// WriteToRequest writes these params to a swagger request
func (o *UploadFileParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Description != nil {

		// form param description
		var frDescription string
		if o.Description != nil {
			frDescription = *o.Description
		}
		fDescription := frDescription
		if fDescription != "" {
			if err := r.SetFormParam("description", fDescription); err != nil {
				return err
			}
		}

	}

	// form file param file
	if err := r.SetFileParam("file", o.File); err != nil {
		return err
	}

	if o.Name != nil {

		// form param name
		var frName string
		if o.Name != nil {
			frName = *o.Name
		}
		fName := frName
		if fName != "" {
			if err := r.SetFormParam("name", fName); err != nil {
				return err
			}
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}