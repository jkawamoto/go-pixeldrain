// create_list.go
//
// Copyright (c) 2018-2021 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package pixeldrain

import (
	"context"
	"strings"

	"github.com/jkawamoto/go-pixeldrain/pkg/pixeldrain/client/list"
)

// CreateList sends a list creation request with the given title, description, and items.
func (pd *Pixeldrain) CreateList(ctx context.Context, title, description string, items []string) (string, error) {
	res, err := pd.cli.List.CreateFileList(
		list.NewCreateFileListParamsWithContext(ctx).WithList(list.CreateFileListBody{
			Title:       &title,
			Description: description,
			Files:       parseListItems(items),
		}),
		pd.authInfoWriter,
	)
	if err != nil {
		return "", NewError(err)
	}

	return res.Payload.ID, nil
}

// parseListItems parses the given list of list specifications and returns a PostListParamsBodyFiles instance.
func parseListItems(values []string) []*list.CreateFileListParamsBodyFilesItems0 {
	res := make([]*list.CreateFileListParamsBodyFilesItems0, len(values))
	for i, v := range values {
		c := strings.SplitN(v, ":", 2)
		item := &list.CreateFileListParamsBodyFilesItems0{
			ID: c[0],
		}
		if len(c) != 1 {
			item.Description = c[1]
		}
		res[i] = item
	}
	return res
}
