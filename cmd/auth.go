// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"log"

	"github.com/rotspace/abudef/core"
	"github.com/rotspace/mtproto"
	"github.com/urfave/cli"
)

var (
	// CmdAuth auth user in t
	CmdAuth = cli.Command{
		Name:   "auth",
		Action: runAuth,
	}
)

func runAuth(*cli.Context) {
	m, err := mtproto.New(mtproto.OptAuthFile(core.DefaultAuthDataFilename))
	if err != nil {
		log.Fatalln(err)
	}

	err = m.Connect()
	if err != nil {
		return
	}

	var (
		phoneNumber string
		code        string
	)

	log.Println("Input phone:")
	fmt.Scanln(&phoneNumber)

	if phoneNumber == "" {
		log.Fatalln("Phone number is empty")
	}
	sentCode, err := m.AuthSendCode(phoneNumber)
	if err != nil {
		log.Fatalf("Err sending code: %s\n", err)
	}

	if !sentCode.Phone_registered {
		log.Fatalf("Phone number isn't registered\n")
	}

	fmt.Printf("Enter code: ")
	fmt.Scanf("%s", &code)
	auth, err := m.AuthSignIn(phoneNumber, code, sentCode.Phone_code_hash)
	if err != nil {
		log.Fatalf("Err autorixation: %s", err)
	}

	userSelf := auth.User.(mtproto.TL_user)
	message := fmt.Sprintf("Signed in: Id %d name <%s @%s %s>\n", userSelf.Id, userSelf.First_name, userSelf.Username, userSelf.Last_name)
	log.Println(message)
}
