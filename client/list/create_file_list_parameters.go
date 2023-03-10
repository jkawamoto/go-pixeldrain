// Code generated by go-swagger; DO NOT EDIT.

package list

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewCreateFileListParams creates a new CreateFileListParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewCreateFileListParams() *CreateFileListParams {
	return &CreateFileListParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewCreateFileListParamsWithTimeout creates a new CreateFileListParams object
// with the ability to set a timeout on a request.
func NewCreateFileListParamsWithTimeout(timeout time.Duration) *CreateFileListParams {
	return &CreateFileListParams{
		timeout: timeout,
	}
}

// NewCreateFileListParamsWithContext creates a new CreateFileListParams object
// with the ability to set a context for a request.
func NewCreateFileListParamsWithContext(ctx context.Context) *CreateFileListParams {
	return &CreateFileListParams{
		Context: ctx,
	}
}

// NewCreateFileListParamsWithHTTPClient creates a new CreateFileListParams object
// with the ability to set a custom HTTPClient for a request.
func NewCreateFileListParamsWithHTTPClient(client *http.Client) *CreateFileListParams {
	return &CreateFileListParams{
		HTTPClient: client,
	}
}

/*
CreateFileListParams contains all the parameters to send to the API endpoint

	for the create file list operation.

	Typically these are written to a http.Request.
*/
type CreateFileListParams struct {

	/* List.

	   POST body should be a JSON object, example below. A list can contain maximally 5000 files. If you try to add more the request will fail.

	*/
	List CreateFileListBody

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the create file list params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *CreateFileListParams) WithDefaults() *CreateFileListParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the create file list params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *CreateFileListParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the create file list params
func (o *CreateFileListParams) WithTimeout(timeout time.Duration) *CreateFileListParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the create file list params
func (o *CreateFileListParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the create file list params
func (o *CreateFileListParams) WithContext(ctx context.Context) *CreateFileListParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the create file list params
func (o *CreateFileListParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the create file list params
func (o *CreateFileListParams) WithHTTPClient(client *http.Client) *CreateFileListParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the create file list params
func (o *CreateFileListParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithList adds the list to the create file list params
func (o *CreateFileListParams) WithList(list CreateFileListBody) *CreateFileListParams {
	o.SetList(list)
	return o
}

// SetList adds the list to the create file list params
func (o *CreateFileListParams) SetList(list CreateFileListBody) {
	o.List = list
}

// WriteToRequest writes these params to a swagger request
func (o *CreateFileListParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error
	if err := r.SetBodyParam(o.List); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
