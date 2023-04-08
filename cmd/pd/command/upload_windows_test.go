//
// upload_windows_test.go
//
// Copyright (c) 2018-2023 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package command

import "testing"

func Test_parseArgument_windows(t *testing.T) {
	tests := []struct {
		arg      string
		wantPath string
		wantName string
	}{
		{
			arg:      "\"C:\\windows\\path_only\"",
			wantPath: "C:\\windows\\path_only",
		},
		{
			arg:      "\"C:\\windows path\":\"quoted name\"",
			wantPath: "C:\\windows path",
			wantName: "quoted name",
		},
		{
			arg:      "C:\\windows\\unquoted_path",
			wantPath: "C:\\windows\\unquoted_path",
		},
		{
			arg:      "C:\\windows unquoted_path:\"quoted name\"",
			wantPath: "C:\\windows unquoted_path",
			wantName: "quoted name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.arg, func(t *testing.T) {
			gotPath, gotName, err := parseArgument(tt.arg)
			if err != nil {
				t.Fatal(err)
			}
			if gotPath != tt.wantPath {
				t.Errorf("parseArgument() gotPath = %v, want %v", gotPath, tt.wantPath)
			}
			if gotName != tt.wantName {
				t.Errorf("parseArgument() gotName = %v, want %v", gotName, tt.wantName)
			}
		})
	}
}
