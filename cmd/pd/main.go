// main.go
//
// Copyright (c) 2018-2019 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = Name
	app.Version = Version
	app.Author = Author
	app.Email = Email
	app.Usage = "Pixeldrain client"

	app.Flags = GlobalFlags
	app.Commands = Commands
	app.CommandNotFound = CommandNotFound
	app.EnableBashCompletion = true
	app.Copyright = fmt.Sprintf("%v <%v>", Author, Email)

	_ = app.Run(os.Args)
}
