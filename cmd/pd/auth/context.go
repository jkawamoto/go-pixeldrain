// context.go
//
// Copyright (c) 2018-2025 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package auth

import (
	"context"

	"github.com/go-openapi/runtime"
)

type ctxMarker struct{}

var (
	ctxMarkerKey = &ctxMarker{}
)

// ToContext attaches the given client to the given context.
func ToContext(ctx context.Context, auth runtime.ClientAuthInfoWriter) context.Context {
	return context.WithValue(ctx, ctxMarkerKey, auth)
}

// Extract takes a client from the given context. If no client is attached to the context, it will return nil.
func Extract(ctx context.Context) runtime.ClientAuthInfoWriter {
	auth, ok := ctx.Value(ctxMarkerKey).(runtime.ClientAuthInfoWriter)
	if !ok {
		return nil
	}
	return auth
}
