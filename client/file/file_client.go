// Code generated by go-swagger; DO NOT EDIT.

package file

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// New creates a new file API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) *Client {
	return &Client{transport: transport, formats: formats}
}

/*
Client for file API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

/*
GetFileIDDownload downloads a file

Same as GET /file/{id}, but with File Transfer HTTP headers. Will trigger a save file dialog when opened in a web browser. Also supports byte range requests.

*/
func (a *Client) GetFileIDDownload(params *GetFileIDDownloadParams) (*GetFileIDDownloadOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetFileIDDownloadParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "GetFileIDDownload",
		Method:             "GET",
		PathPattern:        "/file/{id}/download",
		ProducesMediaTypes: []string{"application/json", "application/octet-stream"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetFileIDDownloadReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*GetFileIDDownloadOK), nil

}

/*
GetFileIDThumbnail gets a thumbnail image representing the file

Returns a PNG thumbnail image representing the file. The thumbnail is always 100*100 px. If the source file is parsable by imagemagick the thumbnail will be generated from the file, if not it will be a generic mime type icon.

*/
func (a *Client) GetFileIDThumbnail(params *GetFileIDThumbnailParams) (*GetFileIDThumbnailOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetFileIDThumbnailParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "GetFileIDThumbnail",
		Method:             "GET",
		PathPattern:        "/file/{id}/thumbnail",
		ProducesMediaTypes: []string{"application/json", "image/png"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetFileIDThumbnailReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*GetFileIDThumbnailOK), nil

}

/*
DeleteFile deletes a file

Deletes a file. Only works when the users owns the file.
*/
func (a *Client) DeleteFile(params *DeleteFileParams) (*DeleteFileOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewDeleteFileParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "deleteFile",
		Method:             "DELETE",
		PathPattern:        "/file/{id}",
		ProducesMediaTypes: []string{"application/json", "text/plain"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &DeleteFileReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*DeleteFileOK), nil

}

/*
DownloadFile downloads a file

Returns the full file associated with the ID. Supports byte range requests.

*/
func (a *Client) DownloadFile(params *DownloadFileParams, writer io.Writer) (*DownloadFileOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewDownloadFileParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "downloadFile",
		Method:             "GET",
		PathPattern:        "/file/{id}",
		ProducesMediaTypes: []string{"application/json", "application/octet-stream"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &DownloadFileReader{formats: a.formats, writer: writer},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*DownloadFileOK), nil

}

/*
GetFileInfo retrieves information of a file

Returns information about one or more files. You can also put a comma separated list of file IDs in the URL and it will return an array of file info, instead of a single object.

*/
func (a *Client) GetFileInfo(params *GetFileInfoParams) (*GetFileInfoOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetFileInfoParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "getFileInfo",
		Method:             "GET",
		PathPattern:        "/file/{id}/info",
		ProducesMediaTypes: []string{"application/json", "text/plain"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetFileInfoReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*GetFileInfoOK), nil

}

/*
UploadFile uploads a file

Upload a file.
*/
func (a *Client) UploadFile(params *UploadFileParams) (*UploadFileCreated, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewUploadFileParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "uploadFile",
		Method:             "POST",
		PathPattern:        "/file",
		ProducesMediaTypes: []string{"application/json", "text/plain"},
		ConsumesMediaTypes: []string{"multipart/form-data"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &UploadFileReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*UploadFileCreated), nil

}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
