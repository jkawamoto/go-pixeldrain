// status.go
//
// Copyright (c) 2018-2021 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package status

const (
	// InvalidCommand is the exit code if the given command is not found.
	InvalidCommand = iota + 1
	// InvalidArgument is the exit code if the given arguments are invalid.
	InvalidArgument
	// APIError is the exit code if an API request fails.
	APIError
)
