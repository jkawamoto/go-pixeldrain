// download.go
//
// Copyright (c) 2018-2025 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package command

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"filippo.io/age"
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/cheggaaa/pb/v3"
	"github.com/go-openapi/swag"
	"github.com/urfave/cli/v2"

	"github.com/jkawamoto/go-pixeldrain"
	"github.com/jkawamoto/go-pixeldrain/client/file"
	"github.com/jkawamoto/go-pixeldrain/client/list"
	"github.com/jkawamoto/go-pixeldrain/cmd/pd/auth"
	"github.com/jkawamoto/go-pixeldrain/cmd/pd/status"
	"github.com/jkawamoto/go-pixeldrain/models"
)

// getID returns the ID from a URL.
func getID(url string) string {
	return url[strings.LastIndex(url, "/")+1:]
}

func downloadURL(ctx *cli.Context, url, dir string, identities []age.Identity) error {
	res, err := pixeldrain.Default.File.GetFileInfo(
		file.NewGetFileInfoParamsWithContext(ctx.Context).WithID(getID(url)),
		auth.Extract(ctx.Context),
	)
	if err != nil {
		return pixeldrain.NewError(err)
	}

	return download(ctx, res.Payload, dir, identities)
}

func download(ctx *cli.Context, info *models.FileInfo, dir string, identities []age.Identity) error {
	var encrypted bool
	if strings.HasSuffix(info.Name, AgeExt) && len(identities) != 0 {
		info.Name = strings.TrimSuffix(info.Name, AgeExt)
		encrypted = true
	}

	filePath := filepath.Join(dir, info.Name)

	// Check if continue flag is set and file already exists with matching size, and skip downloading if so.
	if ctx.Bool(FlagContinue) {
		if stat, err := os.Stat(filePath); err == nil && stat.Size() == info.Size {
			fmt.Fprintf(ctx.App.ErrWriter, "%s (SKIPPED)\n", info.Name)
			return nil
		}
	}

	var w io.WriteCloser
	w, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, w.Close())
	}()
	if encrypted {
		w = Decrypt(w, identities)
	}

	bar := pb.New64(info.Size)
	bar.Set(pb.SIBytesPrefix, true)
	bar.Set("prefix", info.Name+" ")
	bar.SetWriter(ctx.App.ErrWriter)
	bar.Start()
	defer bar.Finish()

	_, _, err = pixeldrain.Default.File.DownloadFile(
		file.NewDownloadFileParamsWithContext(ctx.Context).WithID(swag.StringValue(info.ID)),
		auth.Extract(ctx.Context),
		bar.NewProxyWriter(w),
	)
	if err != nil {
		return pixeldrain.NewError(err)
	}
	return nil
}

func CmdDownload(c *cli.Context) error {
	if c.NArg() == 0 {
		return cli.Exit(fmt.Sprintf("expected at least 1 argument but %d given", c.NArg()), status.InvalidArgument)
	}

	var identities []age.Identity
	if name := c.String(FlagIdentity); name != "" {
		res, err := parseIdentityFile(name)
		if err != nil {
			return cli.Exit(err, status.InvalidArgument)
		}
		identities = res
	}

	dir := c.String(FlagDirectory)
	for _, url := range c.Args().Slice() {
		isList, err := pixeldrain.IsListURL(url)
		if err != nil {
			return cli.Exit(err, status.InvalidArgument)
		}

		// if the given URL doesn't points a list, download it and continue.
		if !isList {
			if err = downloadURL(c, url, dir, identities); err != nil {
				return cli.Exit(err, status.APIError)
			}
			continue
		}

		res, err := pixeldrain.Default.List.GetFileList(
			list.NewGetFileListParamsWithContext(c.Context).WithID(getID(url)),
			auth.Extract(c.Context),
		)
		if err != nil {
			return cli.Exit(pixeldrain.NewError(err), status.APIError)
		}

		all := c.Bool(FlagAll)
		r, ok := c.App.Reader.(terminal.FileReader)
		if !ok {
			all = true
		}
		w, ok := c.App.Writer.(terminal.FileWriter)
		if !ok {
			all = true
		}
		if all {
			for _, f := range res.Payload.Files {
				if err = download(c, f, dir, identities); err != nil {
					return cli.Exit(err, status.APIError)
				}
			}
			continue
		}

		prompt := &survey.MultiSelect{
			Message: "Which files do you want to download:",
			Options: make([]string, len(res.Payload.Files)),
		}
		for i, f := range res.Payload.Files {
			prompt.Options[i] = f.Name
		}

		var names []string
		if err = survey.AskOne(prompt, &names, survey.WithStdio(r, w, c.App.ErrWriter)); err != nil {
			return cli.Exit(err, status.IOError)
		}
		for _, f := range res.Payload.Files {
			if contains(names, f.Name) {
				if err = download(c, f, dir, identities); err != nil {
					return cli.Exit(err, status.APIError)
				}
			}
		}
	}

	return nil
}

func contains(ar []string, v string) bool {
	for _, s := range ar {
		if v == s {
			return true
		}
	}
	return false
}

func parseIdentityFile(name string) (_ []age.Identity, err error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = errors.Join(err, f.Close())
	}()

	return age.ParseIdentities(f)
}
