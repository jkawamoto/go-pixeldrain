// download_test.go
//
// Copyright (c) 2018-2025 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package command

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"io"
	"os"
	"path/filepath"
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

func expectMatchFile(t *testing.T, file1, file2 string) {
	t.Helper()

	data1, err := os.ReadFile(file1)
	if err != nil {
		t.Fatal(err)
	}
	data2, err := os.ReadFile(file2)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data1, data2) {
		t.Errorf("expect matching %v and %v", file1, file2)
	}
}

type GetFileInfoFunc = func(
	*file.GetFileInfoParams,
	runtime.ClientAuthInfoWriter,
	...file.ClientOption,
) (*file.GetFileInfoOK, error)

type DownloadFileFunc = func(
	*file.DownloadFileParams,
	runtime.ClientAuthInfoWriter,
	io.Writer,
	...file.ClientOption,
) (*file.DownloadFileOK, error)

func TestCmdDownload(t *testing.T) {
	apiKey := "test-key"
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatal(err)
	}
	recipient := identity.Recipient()
	dir := t.TempDir()

	getFileInfo := func(ctx context.Context, encrypted bool) GetFileInfoFunc {
		return func(
			params *file.GetFileInfoParams,
			authInfo runtime.ClientAuthInfoWriter,
			opts ...file.ClientOption,
		) (*file.GetFileInfoOK, error) {
			if params.Context != ctx {
				t.Errorf("expect %v, got %v", ctx, params.Context)
			}
			testutil.ExpectAuthInfoWritesAPIKey(t, authInfo, apiKey)

			info, err := os.Stat(params.ID)
			if err != nil {
				t.Fatal(err)
			}
			ext := ".download"
			if encrypted {
				ext = ext + ".age"
			}
			return &file.GetFileInfoOK{
				Payload: &models.FileInfo{
					ID:   swag.String(params.ID),
					Name: params.ID + ext,
					Size: info.Size(),
				},
			}, nil

		}
	}
	downloadFile := func(ctx context.Context, recipient age.Recipient) DownloadFileFunc {
		return func(
			params *file.DownloadFileParams,
			authInfo runtime.ClientAuthInfoWriter,
			writer io.Writer,
			opts ...file.ClientOption,
		) (*file.DownloadFileOK, error) {
			if params.Context != ctx {
				t.Errorf("expect %v, got %v", ctx, params.Context)
			}
			testutil.ExpectAuthInfoWritesAPIKey(t, authInfo, apiKey)

			f, err := os.Open(params.ID)
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				if err := f.Close(); err != nil {
					t.Fatal(err)
				}
			}()
			if recipient != nil {
				w, err := age.Encrypt(writer, recipient)
				if err != nil {
					t.Fatal(err)
				}
				defer func() {
					if err := w.Close(); err != nil {
						t.Fatal(err)
					}
				}()
				writer = w
			}
			if _, err = io.Copy(writer, f); err != nil {
				t.Fatal(err)
			}

			return &file.DownloadFileOK{
				Payload: writer,
			}, nil
		}
	}

	cases := []struct {
		name   string
		init   func(*testing.T, context.Context, *mock.MockClientService)
		args   []string
		expect []string
		exit   int
	}{
		{
			name: "not enough arguments",
			exit: status.InvalidArgument,
		},
		{
			name: "download one file",
			init: func(t *testing.T, ctx context.Context, m *mock.MockClientService) {
				t.Helper()

				m.EXPECT().GetFileInfo(gomock.Any(), gomock.Any()).DoAndReturn(getFileInfo(ctx, false))
				m.EXPECT().DownloadFile(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(downloadFile(ctx, nil))
			},
			args:   []string{pixeldrain.DownloadURL("doc.go")},
			expect: []string{"doc.go"},
		},
		{
			name: "download multiple files",
			init: func(t *testing.T, ctx context.Context, m *mock.MockClientService) {
				t.Helper()

				m.EXPECT().
					GetFileInfo(gomock.Any(), gomock.Any()).
					DoAndReturn(getFileInfo(ctx, false)).
					Times(2)
				m.EXPECT().
					DownloadFile(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(downloadFile(ctx, nil)).
					Times(2)
			},
			args:   []string{pixeldrain.DownloadURL("doc.go"), pixeldrain.DownloadURL("download.go")},
			expect: []string{"doc.go", "download.go"},
		},
		{
			name: "download one list",
			init: func(t *testing.T, ctx context.Context, m *mock.MockClientService) {
				t.Helper()

				m.EXPECT().
					GetFileList(gomock.Any(), gomock.Any()).
					DoAndReturn(func(
						params *list.GetFileListParams,
						authInfo runtime.ClientAuthInfoWriter,
						opts ...list.ClientOption,
					) (*list.GetFileListOK, error) {
						if params.Context != ctx {
							t.Errorf("expect %v, got %v", ctx, params.Context)
						}
						testutil.ExpectAuthInfoWritesAPIKey(t, authInfo, apiKey)

						if params.ID != "abc" {
							t.Errorf("expect %v, got %v", "abc", params.ID)
						}

						return &list.GetFileListOK{
							Payload: &models.GetFileListResponse{
								Files: []*models.FileInfo{
									{
										ID:   swag.String("doc.go"),
										Name: "doc.go.download",
									},
									{
										ID:   swag.String("download.go"),
										Name: "download.go.download",
									},
								},
								ID:      params.ID,
								Success: true,
							},
						}, nil
					})

				m.EXPECT().
					DownloadFile(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(downloadFile(ctx, nil)).
					Times(2)
			},
			args:   []string{pixeldrain.ListURL("abc")},
			expect: []string{"doc.go", "download.go"},
		},
		{
			name: "download one encrypted file",
			init: func(t *testing.T, ctx context.Context, m *mock.MockClientService) {
				t.Helper()

				err := os.WriteFile(filepath.Join(dir, "key.txt"), []byte(identity.String()), 0600)
				if err != nil {
					t.Fatal(err)
				}

				m.EXPECT().GetFileInfo(gomock.Any(), gomock.Any()).DoAndReturn(getFileInfo(ctx, true))
				m.EXPECT().
					DownloadFile(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(downloadFile(ctx, recipient))
			},
			args:   []string{"--identity", filepath.Join(dir, "key.txt"), pixeldrain.DownloadURL("doc.go")},
			expect: []string{"doc.go"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			dir := t.TempDir()

			flagSet := flag.NewFlagSet("download", flag.PanicOnError)
			flagSet.String(FlagDirectory, dir, "")
			flagSet.Bool(FlagAll, true, "")
			flagSet.String(FlagIdentity, "", "")
			err := flagSet.Parse(tc.args)
			if err != nil {
				t.Fatal(err)
			}

			c := cli.NewContext(cli.NewApp(), flagSet, nil)
			c.Context = auth.ToContext(c.Context, client.BasicAuth("", apiKey))

			m := mock.NewMockClientService(ctrl)
			if tc.init != nil {
				tc.init(t, c.Context, m)
			}
			RegisterMock(t, m)

			err = CmdDownload(c)
			if err != nil || tc.exit != 0 {
				var e cli.ExitCoder
				if !errors.As(err, &e) {
					t.Errorf("expect an ExitCoder, got %v", err)
				} else if e.ExitCode() != tc.exit {
					t.Errorf("expect %v, got %v", tc.exit, e.ExitCode())
				}
			}

			for _, name := range tc.expect {
				expectMatchFile(t, name, filepath.Join(dir, name+".download"))
			}
		})
	}
}
