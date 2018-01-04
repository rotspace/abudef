// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/rotspace/abudef/core"
	"github.com/urfave/cli"
)

var (
	// CmdBot is cli command
	CmdBot = cli.Command{
		Name:   "bot",
		Action: runBot,
	}
)

func runBot(*cli.Context) {
	log.SetFlags(log.Llongfile | log.LstdFlags)

	var (
		envFile = "/home/admin/web/redcastle.info/.env"
		key     = "TELEGRAM_BOT_TOKEN"
		cbURL   = "https://redcastle.info/setWebhook"
	)
	eng, err := core.New(envFile, key, cbURL)
	if err != nil {
		log.Fatalln(err)
	}

	mtp, err := core.NewDefaultMTPRotoWrapper()
	if err != nil {
		log.Fatalln(err)
	}
	eng.Tg = mtp

	store, err := core.NewStore()
	if err != nil {
		log.Fatalln(err)
	}
	eng.Store = store

	eng.OnMessage(onMessage)

	eng.Run()
}

func onMessage(currentUsername, unused string) (string, error) {
	bts, err := ioutil.ReadFile("text.txt")
	if err != nil {
		return "@" + currentUsername, nil
	}
	return fmt.Sprintf(string(bts), currentUsername), nil
}
