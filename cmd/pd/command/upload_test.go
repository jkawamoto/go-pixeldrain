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
	"flag"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"filippo.io/age"
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
	"github.com/jkawamoto/go-pixeldrain/cmd/pd/internal/testutil"
	"github.com/jkawamoto/go-pixeldrain/cmd/pd/status"
	"github.com/jkawamoto/go-pixeldrain/models"
)

func TestUpload(t *testing.T) {
	apiKey := "test-key"
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatal(err)
	}
	recipient := identity.Recipient()
	dir := t.TempDir()

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
						testutil.ExpectAuthInfoWritesAPIKey(t, authInfo, apiKey)

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
						testutil.ExpectAuthInfoWritesAPIKey(t, authInfo, apiKey)

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
						testutil.ExpectAuthInfoWritesAPIKey(t, authInfo, apiKey)

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
						testutil.ExpectAuthInfoWritesAPIKey(t, authInfo, apiKey)

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
						testutil.ExpectAuthInfoWritesAPIKey(t, authInfo, apiKey)

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
						testutil.ExpectAuthInfoWritesAPIKey(t, authInfo, apiKey)

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
						testutil.ExpectAuthInfoWritesAPIKey(t, authInfo, apiKey)

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
		{
			name: "encrypt with a recipient's public key",
			args: []string{"--recipient", recipient.String(), "doc.go"},
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
						testutil.ExpectAuthInfoWritesAPIKey(t, authInfo, apiKey)

						r, err := age.Decrypt(params.File, identity)
						if err != nil {
							t.Fatal(err)
						}
						data, err := io.ReadAll(r)
						if err != nil {
							t.Fatal(err)
						}
						if err = params.File.Close(); err != nil {
							t.Fatal(err)
						}

						expect, err := os.ReadFile(strings.TrimSuffix(params.File.Name(), AgeExt))
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
					})
			},
			expect: pixeldrain.DownloadURL("doc.go.age") + "\n",
		},
		{
			name: "encrypt with a recipient file",
			args: []string{"--recipient-file", filepath.Join(dir, "key.txt"), "doc.go"},
			init: func(t *testing.T, ctx context.Context, m *mock.MockClientService) {
				t.Helper()

				err := os.WriteFile(filepath.Join(dir, "key.txt"), []byte(recipient.String()), 0600)
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
						testutil.ExpectAuthInfoWritesAPIKey(t, authInfo, apiKey)

						r, err := age.Decrypt(params.File, identity)
						if err != nil {
							t.Fatal(err)
						}
						data, err := io.ReadAll(r)
						if err != nil {
							t.Fatal(err)
						}
						if err = params.File.Close(); err != nil {
							t.Fatal(err)
						}

						expect, err := os.ReadFile(strings.TrimSuffix(params.File.Name(), AgeExt))
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
					})
			},
			expect: pixeldrain.DownloadURL("doc.go.age") + "\n",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			flagSet := flag.NewFlagSet("upload", flag.PanicOnError)
			flagSet.String(FlagAlbumName, "", "")
			flagSet.Var(&cli.StringSlice{}, FlagRecipient, "")
			flagSet.String(FlagRecipientFile, "", "")
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

func Test_parseArgument(t *testing.T) {
	tests := []struct {
		arg      string
		wantPath string
		wantName string
	}{
		{
			arg:      "path_only",
			wantPath: "path_only",
		},
		{
			arg:      "simple:case",
			wantPath: "simple",
			wantName: "case",
		},
		{
			arg:      "\"quoted path only\"",
			wantPath: "quoted path only",
			wantName: "",
		},
		{
			arg:      "\"quoted path only\":\"quoted name\"",
			wantPath: "quoted path only",
			wantName: "quoted name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.arg, func(t *testing.T) {
			gotPath, gotName, err := parseArgument(tt.arg)
			if err != nil {
				t.Fatal(err)
			}
			if gotPath != tt.wantPath {
				t.Errorf("parseArgument() gotPath = %v, want %v", gotPath, tt.wantPath)
			}
			if gotName != tt.wantName {
				t.Errorf("parseArgument() gotName = %v, want %v", gotName, tt.wantName)
			}
		})
	}
}
