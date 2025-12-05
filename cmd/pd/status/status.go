// status.go
//
// Copyright (c) 2018-2025 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

// Package status defines status codes pd command returns.
package status

const (
	// InvalidCommand is the exit code if the given command is not found.
	InvalidCommand = iota + 1
	// InvalidArgument is the exit code if the given arguments are invalid.
	InvalidArgument
	// APIError is the exit code if an API request fails.
	APIError
	// IOError is the exit code if an IO error happens.
	IOError
)
