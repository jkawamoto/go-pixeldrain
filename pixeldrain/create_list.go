// create_list.go
//
// Copyright (c) 2018 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package pixeldrain

import (
	"context"
	"strings"

	"github.com/jkawamoto/go-pixeldrain/client"
	"github.com/jkawamoto/go-pixeldrain/client/list"
	"github.com/jkawamoto/go-pixeldrain/models"
)

// CreateList sends a list creation request via the given client with the given title, description, and items.
func CreateList(ctx context.Context, cli *client.PixelDrain, title, description string, items []string) (id string, err error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if cli == nil {
		cli = client.Default
	}

	res, err := cli.List.PostList(list.NewPostListParamsWithContext(ctx).WithList(
		&models.PostListParamsBody{
			Title:       &title,
			Description: description,
			Files:       ParseListItems(items),
		},
	))
	if err != nil {
		return
	}

	id = res.Payload.ID
	return

}

// ParseListItems parses the given list of list specifications and returns a PostListParamsBodyFiles instance.
func ParseListItems(values []string) (res models.PostListParamsBodyFiles) {

	for _, v := range values {
		c := strings.SplitN(v, ":", 2)
		item := &models.PostListParamsBodyFilesItems{
			ID: c[0],
		}
		if len(c) != 1 {
			item.Description = c[1]
		}
		res = append(res, item)
	}
	return
}
