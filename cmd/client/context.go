// context.go
//
// Copyright (c) 2018-2023 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package client

import (
	"context"
)

type ctxMarker struct{}

var (
	ctxMarkerKey = &ctxMarker{}
)

// ToContext attaches the given client to the given context.
func ToContext(ctx context.Context, cli Client) context.Context {
	return context.WithValue(ctx, ctxMarkerKey, cli)
}

// Extract takes a client from the given context. If no client is attached to the context, it will return nil.
func Extract(ctx context.Context) Client {
	cli, ok := ctx.Value(ctxMarkerKey).(Client)
	if !ok {
		return nil
	}
	return cli
}