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
func (pd *Pixeldrain) CreateList(ctx context.Context, title, description string, items []string) (id string, err error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	res, err := pd.Client.List.CreateFileList(list.NewCreateFileListParamsWithContext(ctx).WithList(
		list.CreateFileListBody{
			Title:       &title,
			Description: description,
			Files:       parseListItems(items),
		},
	))
	if err != nil {
		return
	}

	id = res.Payload.ID
	return

}

// parseListItems parses the given list of list specifications and returns a PostListParamsBodyFiles instance.
func parseListItems(values []string) (res []*list.FilesItems0) {

	for _, v := range values {
		c := strings.SplitN(v, ":", 2)
		item := &list.FilesItems0{
			ID: c[0],
		}
		if len(c) != 1 {
			item.Description = c[1]
		}
		res = append(res, item)
	}
	return
}
