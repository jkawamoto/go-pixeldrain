// main_test.go
//
// Copyright (c) 2018-2025 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package main

import (
	"flag"
	"testing"

	"github.com/urfave/cli/v2"

	"github.com/jkawamoto/go-pixeldrain/cmd/pd/auth"
	"github.com/jkawamoto/go-pixeldrain/cmd/pd/internal/testutil"
)

func Test_initApp(t *testing.T) {
	apiKey := "test-api-key"

	flagSet := flag.NewFlagSet("", flag.PanicOnError)
	flagSet.String(FlagAPIKey, apiKey, "")
	err := flagSet.Parse([]string{})
	if err != nil {
		t.Fatal(err)
	}

	c := cli.NewContext(cli.NewApp(), flagSet, nil)

	app := initApp()
	if err = app.Before(c); err != nil {
		t.Fatal(err)
	}

	testutil.ExpectAuthInfoWritesAPIKey(t, auth.Extract(c.Context), apiKey)
}
