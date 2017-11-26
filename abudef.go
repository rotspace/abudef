// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"os"

	"github.com/rotspace/abudef/cmd"
	"github.com/urfave/cli"
)

// Version is current app version
const Version = "0.0.1"

func main() {
	app := &cli.App{
		Version: Version,
		Commands: cli.Commands{
			cmd.CmdBot,
			cmd.CmdAuth,
		},
	}
	app.Run(os.Args)
}
