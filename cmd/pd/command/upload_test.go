// upload_test.go
//
// Copyright (c) 2018-2023 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package command

import (
	"bytes"
	"context"
	"encoding/base64"
	"flag"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/client"
	"github.com/go-openapi/swag"
	"github.com/golang/mock/gomock"
	"github.com/urfave/cli/v2"

	"github.com/jkawamoto/go-pixeldrain"
	"github.com/jkawamoto/go-pixeldrain/client/file"
	"github.com/jkawamoto/go-pixeldrain/client/list"
	"github.com/jkawamoto/go-pixeldrain/cmd/pd/auth"
	"github.com/jkawamoto/go-pixeldrain/cmd/pd/command/mock"
	"github.com/jkawamoto/go-pixeldrain/cmd/pd/status"
	"github.com/jkawamoto/go-pixeldrain/models"
)

func checkAPIKey(t *testing.T, authInfo runtime.ClientAuthInfoWriter, apiKey string) {
	t.Helper()

	req := &runtime.TestClientRequest{}
	if err := authInfo.AuthenticateRequest(req, nil); err != nil {
		t.Fatal(err)
	}
	s, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(req.Headers.Get(runtime.HeaderAuthorization), "Basic "))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(s), apiKey) {
		t.Errorf("expect api key, got %v", string(s))
	}
}

func TestUpload(t *testing.T) {
	apiKey := "test-key"

	cases := []struct {
		name   string
		init   func(*testing.T, context.Context, *mock.MockClientService)
		args   []string
		stdin  io.Reader
		expect string
		exit   int
	}{
		{
			name: "not enough arguments",
			exit: status.InvalidArgument,
		},
		{
			name: "one file w/o renaming",
			args: []string{"doc.go"},
			init: func(t *testing.T, ctx context.Context, m *mock.MockClientService) {
				t.Helper()
				expect, err := os.ReadFile("doc.go")
				if err != nil {
					t.Fatal(err)
				}

				m.EXPECT().
					UploadFile(gomock.Any(), gomock.Any()).
					DoAndReturn(func(
						params *file.UploadFileParams,
						authInfo runtime.ClientAuthInfoWriter,
						opts ...file.ClientOption,
					) (*file.UploadFileCreated, error) {
						if params.Context != ctx {
							t.Errorf("expect %v, got %v", ctx, params.Context)
						}
						checkAPIKey(t, authInfo, apiKey)

						if name := params.File.Name(); name != "doc.go" {
							t.Errorf("expect %v, got %v", "doc.go", name)
						}
						data, err := io.ReadAll(params.File)
						if err != nil {
							t.Fatal(err)
						}
						if err = params.File.Close(); err != nil {
							t.Error(err)
						}
						if !bytes.Equal(expect, data) {
							t.Errorf("expect %v, got %v", expect, data)
						}
						return &file.UploadFileCreated{
							Payload: &models.SuccessResponse{
								ID:      swag.String("123"),
								Success: swag.Bool(true),
							},
						}, nil
					})
			},
			expect: pixeldrain.DownloadURL("123") + "\n",
		},
		{
			name: "one file w/ renaming",
			args: []string{"doc.go:manual"},
			init: func(t *testing.T, ctx context.Context, m *mock.MockClientService) {
				t.Helper()
				expect, err := os.ReadFile("doc.go")
				if err != nil {
					t.Fatal(err)
				}

				m.EXPECT().
					UploadFile(gomock.Any(), gomock.Any()).
					DoAndReturn(func(
						params *file.UploadFileParams,
						authInfo runtime.ClientAuthInfoWriter,
						opts ...file.ClientOption,
					) (*file.UploadFileCreated, error) {
						if params.Context != ctx {
							t.Errorf("expect %v, got %v", ctx, params.Context)
						}
						checkAPIKey(t, authInfo, apiKey)

						if name := params.File.Name(); name != "manual" {
							t.Errorf("expect %v, got %v", "manual", name)
						}
						data, err := io.ReadAll(params.File)
						if err != nil {
							t.Fatal(err)
						}
						if err = params.File.Close(); err != nil {
							t.Fatal(err)
						}

						if !bytes.Equal(expect, data) {
							t.Errorf("expect %v, got %v", expect, data)
						}
						return &file.UploadFileCreated{
							Payload: &models.SuccessResponse{
								ID:      swag.String("123"),
								Success: swag.Bool(true),
							},
						}, nil
					})
			},
			expect: pixeldrain.DownloadURL("123") + "\n",
		},
		{
			name:  "read from stdin",
			args:  []string{"--", "-:manual"},
			stdin: bytes.NewReader([]byte("test data")),
			init: func(t *testing.T, ctx context.Context, m *mock.MockClientService) {
				t.Helper()
				expect := []byte("test data")

				m.EXPECT().
					UploadFile(gomock.Any(), gomock.Any()).
					DoAndReturn(func(
						params *file.UploadFileParams,
						authInfo runtime.ClientAuthInfoWriter,
						opts ...file.ClientOption,
					) (*file.UploadFileCreated, error) {
						if params.Context != ctx {
							t.Errorf("expect %v, got %v", ctx, params.Context)
						}
						checkAPIKey(t, authInfo, apiKey)

						if name := params.File.Name(); name != "manual" {
							t.Errorf("expect %v, got %v", "manual", name)
						}
						data, err := io.ReadAll(params.File)
						if err != nil {
							t.Fatal(err)
						}
						if err = params.File.Close(); err != nil {
							t.Fatal(err)
						}
						if !bytes.Equal(expect, data) {
							t.Errorf("expect %v, got %v", expect, data)
						}
						return &file.UploadFileCreated{
							Payload: &models.SuccessResponse{
								ID:      swag.String("123"),
								Success: swag.Bool(true),
							},
						}, nil
					})
			},
			expect: pixeldrain.DownloadURL("123") + "\n",
		},

		{
			name: "create list w/ name",
			args: []string{"-album", "list", "doc.go", "upload.go"},
			init: func(t *testing.T, ctx context.Context, m *mock.MockClientService) {
				t.Helper()

				m.EXPECT().
					UploadFile(gomock.Any(), gomock.Any()).
					DoAndReturn(func(
						params *file.UploadFileParams,
						authInfo runtime.ClientAuthInfoWriter,
						opts ...file.ClientOption,
					) (*file.UploadFileCreated, error) {
						if params.Context != ctx {
							t.Errorf("expect %v, got %v", ctx, params.Context)
						}
						checkAPIKey(t, authInfo, apiKey)

						data, err := io.ReadAll(params.File)
						if err != nil {
							t.Fatal(err)
						}
						if err = params.File.Close(); err != nil {
							t.Fatal(err)
						}

						expect, err := os.ReadFile(params.File.Name())
						if err != nil {
							t.Fatal(err)
						}
						if !bytes.Equal(expect, data) {
							t.Errorf("expect %v, got %v", expect, data)
						}
						return &file.UploadFileCreated{
							Payload: &models.SuccessResponse{
								ID:      swag.String(params.File.Name()),
								Success: swag.Bool(true),
							},
						}, nil
					}).Times(2)
				m.EXPECT().
					CreateFileList(gomock.Any(), gomock.Any()).
					DoAndReturn(func(
						params *list.CreateFileListParams,
						authInfo runtime.ClientAuthInfoWriter,
						opts ...list.ClientOption,
					) (*list.CreateFileListCreated, error) {
						if params.Context != ctx {
							t.Errorf("expect %v, got %v", ctx, params.Context)
						}
						checkAPIKey(t, authInfo, apiKey)

						if len(params.List.Files) != 2 {
							t.Fatalf("expect %v, got %v", 2, len(params.List.Files))
						}
						if id := swag.StringValue(params.List.Files[0].ID); id != "doc.go" {
							t.Errorf("expect %v, got %v", "doc.go", id)
						}
						if id := swag.StringValue(params.List.Files[1].ID); id != "upload.go" {
							t.Errorf("expect %v, got %v", "upload.go", id)
						}
						if name := swag.StringValue(params.List.Title); name != "list" {
							t.Errorf("expect %v, got %v", "list", name)
						}

						return &list.CreateFileListCreated{
							Payload: &models.SuccessResponse{
								ID:      swag.String("abc"),
								Success: swag.Bool(true),
							},
						}, nil
					})
			},
			expect: pixeldrain.ListURL("abc") + "\n",
		},
		{
			name: "create list w/o name",
			args: []string{"doc.go", "upload.go"},
			init: func(t *testing.T, ctx context.Context, m *mock.MockClientService) {
				t.Helper()

				m.EXPECT().
					UploadFile(gomock.Any(), gomock.Any()).
					DoAndReturn(func(
						params *file.UploadFileParams,
						authInfo runtime.ClientAuthInfoWriter,
						opts ...file.ClientOption,
					) (*file.UploadFileCreated, error) {
						if params.Context != ctx {
							t.Errorf("expect %v, got %v", ctx, params.Context)
						}
						checkAPIKey(t, authInfo, apiKey)

						data, err := io.ReadAll(params.File)
						if err != nil {
							t.Fatal(err)
						}
						if err = params.File.Close(); err != nil {
							t.Fatal(err)
						}

						expect, err := os.ReadFile(params.File.Name())
						if err != nil {
							t.Fatal(err)
						}
						if !bytes.Equal(expect, data) {
							t.Errorf("expect %v, got %v", expect, data)
						}
						return &file.UploadFileCreated{
							Payload: &models.SuccessResponse{
								ID:      swag.String(params.File.Name()),
								Success: swag.Bool(true),
							},
						}, nil
					}).Times(2)
				m.EXPECT().
					CreateFileList(gomock.Any(), gomock.Any()).
					DoAndReturn(func(
						params *list.CreateFileListParams,
						authInfo runtime.ClientAuthInfoWriter,
						opts ...list.ClientOption,
					) (*list.CreateFileListCreated, error) {
						if params.Context != ctx {
							t.Errorf("expect %v, got %v", ctx, params.Context)
						}
						checkAPIKey(t, authInfo, apiKey)

						if name := swag.StringValue(params.List.Title); !strings.HasPrefix(name, "album-") {
							t.Errorf("expect having prefix album-, got %v", name)
						}
						return &list.CreateFileListCreated{
							Payload: &models.SuccessResponse{
								ID:      swag.String("abc"),
								Success: swag.Bool(true),
							},
						}, nil
					})
			},
			expect: pixeldrain.ListURL("abc") + "\n",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			flagSet := flag.NewFlagSet("upload", flag.PanicOnError)
			flagSet.String(FlagAlbumName, "", "")
			err := flagSet.Parse(tc.args)
			if err != nil {
				t.Fatal(err)
			}

			buf := bytes.NewBuffer(nil)

			c := cli.NewContext(cli.NewApp(), flagSet, nil)
			c.App.Reader = tc.stdin
			c.App.Writer = buf
			c.Context = auth.ToContext(c.Context, client.BasicAuth("", apiKey))

			m := mock.NewMockClientService(ctrl)
			if tc.init != nil {
				tc.init(t, c.Context, m)
			}
			RegisterMock(t, m)

			err = CmdUpload(c)
			if err != nil || tc.exit != 0 {
				if e, ok := err.(cli.ExitCoder); !ok {
					t.Errorf("expect an ExitCoder, got %v", err)
				} else if e.ExitCode() != tc.exit {
					t.Errorf("expect %v, got %v", tc.exit, e.ExitCode())
				}
			}

			if res := buf.String(); res != tc.expect {
				t.Errorf("expect %q, got %q", tc.expect, res)
			}
		})
	}
}
