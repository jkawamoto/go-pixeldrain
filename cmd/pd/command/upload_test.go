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
	"errors"
	"flag"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/urfave/cli/v2"

	"github.com/jkawamoto/go-pixeldrain"
	"github.com/jkawamoto/go-pixeldrain/cmd/client"
	"github.com/jkawamoto/go-pixeldrain/cmd/client/mock"
)

func TestUpload(t *testing.T) {
	cases := []struct {
		name   string
		init   func(*testing.T, context.Context, *mock.MockClient)
		args   []string
		stdin  io.Reader
		expect string
		err    error
	}{
		{
			name: "not enough arguments",
			err:  ErrNotEnoughArguments,
		},
		{
			name: "one file w/o renaming",
			args: []string{"doc.go"},
			init: func(t *testing.T, ctx context.Context, m *mock.MockClient) {
				t.Helper()
				expect, err := os.ReadFile("doc.go")
				if err != nil {
					t.Fatal(err)
				}

				m.EXPECT().
					Upload(ctx, gomock.Any()).
					DoAndReturn(func(_ context.Context, f pixeldrain.File) (string, error) {
						if f.Name() != "doc.go" {
							t.Errorf("expect %v, got %v", "doc.go", f.Name())
						}
						data, err := io.ReadAll(f)
						if err != nil {
							t.Fatal(err)
						}
						if !bytes.Equal(expect, data) {
							t.Errorf("expect %v, got %v", expect, data)
						}
						return "123", nil
					})
				m.EXPECT().DownloadURL("123").Return("https://example.com/123")
			},
			expect: "https://example.com/123\n",
		},
		{
			name: "one file w/ renaming",
			args: []string{"doc.go:manual"},
			init: func(t *testing.T, ctx context.Context, m *mock.MockClient) {
				t.Helper()
				expect, err := os.ReadFile("doc.go")
				if err != nil {
					t.Fatal(err)
				}

				m.EXPECT().
					Upload(ctx, gomock.Any()).
					DoAndReturn(func(_ context.Context, f pixeldrain.File) (string, error) {
						if f.Name() != "manual" {
							t.Errorf("expect %v, got %v", "manual", f.Name())
						}
						data, err := io.ReadAll(f)
						if err != nil {
							t.Fatal(err)
						}
						if !bytes.Equal(expect, data) {
							t.Errorf("expect %v, got %v", expect, data)
						}
						return "123", nil
					})
				m.EXPECT().DownloadURL("123").Return("https://example.com/123")
			},
			expect: "https://example.com/123\n",
		},
		{
			name:  "read from stdin",
			args:  []string{"--", "-:manual"},
			stdin: bytes.NewReader([]byte("test data")),
			init: func(t *testing.T, ctx context.Context, m *mock.MockClient) {
				t.Helper()
				expect := []byte("test data")

				m.EXPECT().
					Upload(ctx, gomock.Any()).
					DoAndReturn(func(_ context.Context, f pixeldrain.File) (string, error) {
						if f.Name() != "manual" {
							t.Errorf("expect %v, got %v", "manual", f.Name())
						}
						data, err := io.ReadAll(f)
						if err != nil {
							t.Fatal(err)
						}
						if !bytes.Equal(expect, data) {
							t.Errorf("expect %v, got %v", expect, data)
						}
						return "123", nil
					})
				m.EXPECT().DownloadURL("123").Return("https://example.com/123")
			},
			expect: "https://example.com/123\n",
		},
		{
			name: "create list w/ name",
			args: []string{"-album", "list", "doc.go", "upload.go"},
			init: func(t *testing.T, ctx context.Context, m *mock.MockClient) {
				t.Helper()

				m.EXPECT().
					Upload(ctx, gomock.Any()).
					DoAndReturn(func(_ context.Context, f pixeldrain.File) (string, error) {
						data, err := io.ReadAll(f)
						if err != nil {
							t.Fatal(err)
						}
						expect, err := os.ReadFile(f.Name())
						if err != nil {
							t.Fatal(err)
						}
						if !bytes.Equal(expect, data) {
							t.Errorf("expect %v, got %v", expect, data)
						}
						return f.Name(), nil
					}).Times(2)
				m.EXPECT().CreateList(ctx, "list", []string{"doc.go", "upload.go"}).Return("abc", nil)

				m.EXPECT().ListURL("abc").Return("https://example.com/abc")
			},
			expect: "https://example.com/abc\n",
		},
		{
			name: "create list w/o name",
			args: []string{"doc.go", "upload.go"},
			init: func(t *testing.T, ctx context.Context, m *mock.MockClient) {
				t.Helper()

				m.EXPECT().
					Upload(ctx, gomock.Any()).
					DoAndReturn(func(_ context.Context, f pixeldrain.File) (string, error) {
						data, err := io.ReadAll(f)
						if err != nil {
							t.Fatal(err)
						}
						expect, err := os.ReadFile(f.Name())
						if err != nil {
							t.Fatal(err)
						}
						if !bytes.Equal(expect, data) {
							t.Errorf("expect %v, got %v", expect, data)
						}
						return f.Name(), nil
					}).Times(2)
				m.EXPECT().
					CreateList(ctx, gomock.Any(), []string{"doc.go", "upload.go"}).
					DoAndReturn(func(_ context.Context, name string, _ []string) (string, error) {
						if !strings.HasPrefix(name, "album-") {
							t.Errorf("expect having prefix album-, got %v", name)
						}
						return "abc", nil
					})

				m.EXPECT().ListURL("abc").Return("https://example.com/abc")
			},
			expect: "https://example.com/abc\n",
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

			m := mock.NewMockClient(ctrl)
			c := cli.NewContext(cli.NewApp(), flagSet, nil)
			c.App.Reader = tc.stdin
			c.App.Writer = buf
			c.Context = client.ToContext(c.Context, m)
			if tc.init != nil {
				tc.init(t, c.Context, m)
			}

			err = CmdUpload(c)
			if !errors.Is(err, tc.err) {
				t.Errorf("expect %v, got %v", tc.err, err)
			}
			if res := buf.String(); res != tc.expect {
				t.Errorf("expect %v, got %v", tc.expect, res)
			}
		})
	}
}
