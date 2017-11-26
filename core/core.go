// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/fatih/color"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Core main structs
type Core struct {
	// Tg is telegram api client
	Tg MTProtoManager

	Store Store

	CurrentToken string

	// CurrentBotUsername is just cache field used for answer in incoming messages
	CurrentBotUsername string

	envFile   string
	valueKey  string
	notifyURL string

	onMessage func(currentBotUsername string, msg string) (answer string, err error)
}

// New returns new instance of core
func New(envFile string, valueKey, notifyURL string) (_ *Core, err error) {
	c := new(Core)

	c.envFile = envFile
	c.valueKey = valueKey
	c.notifyURL = notifyURL

	currentToken, err := getFileValue(c.envFile, valueKey)
	if err != nil {
		return
	}

	username, err := c.CheckToken(currentToken)
	c.CurrentBotUsername = username
	if err != nil {
		err = c.GetNewToken()
		if err != nil {
			return
		}
	}

	return c, nil
}

// OnMessage set message handlers
func (c *Core) OnMessage(fn func(currentBotUsername, msg string) (answer string, err error)) {
	c.onMessage = fn
}

// Run check new messages and answer it
func (c *Core) Run() {
	for {
		_, err := c.CheckToken(c.CurrentToken)
		if err != nil {
			err = c.GetNewToken()
			if err != nil {
				color.Red("Error getting new valid token %s", err)
			}
		}

		// handle new messages
		if c.onMessage != nil {
			dialogs, err := c.Tg.GetUnreadedDialogs()
			if err != nil {
				color.Red("Error getting unread dialogs %s", err)
			}

			message, err := c.onMessage(c.CurrentBotUsername, "")
			if err != nil {
				time.Sleep(1 * time.Minute)
				continue
			}

			for _, dia := range dialogs {
				err = c.Tg.SendMessage(dia, message)
				if err != nil {
					color.Red("Error sending message %s", err)
				}
				time.Sleep(1 * time.Second)
			}
		}

		time.Sleep(1 * time.Minute)
	}
}

// CheckToken checks telegram token, returns error if token invalid
func (c *Core) CheckToken(tok string) (username string, err error) {
	t, err := tgbotapi.NewBotAPI(tok)
	if err != nil {
		return
	}
	bot, err := t.GetMe()
	if err != nil {
		return
	}
	username = bot.UserName
	return
}

// SetValidToken set valid token to current bot
func (c *Core) SetValidToken(token, username string) (err error) {
	c.CurrentBotUsername = username
	c.CurrentToken = token

	err = replaceFileValue(c.envFile, c.valueKey, token)
	if err != nil {
		return
	}

	err = notifyHTTP(c.notifyURL)

	return
}

func notifyHTTP(uri string) (err error) {
	resp, err := http.Get(uri)
	if err != nil {
		return
	}
	resp.Body.Close()

	// TODO: check is url == / (index)
	//if resp.Request.URL.String() !=

	return
}

func getFileValue(filename, key string) (value string, err error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	reg := fmt.Sprintf(`(%s ?= ?(\S*))`, key)
	re := regexp.MustCompile(reg)
	res := re.FindSubmatch(data)
	if len(res) != 3 {
		err = fmt.Errorf("Value not found")
		return
	}
	return string(res[2]), nil
}

// replace value in bash varianle definition file
func replaceFileValue(filename, key, newValue string) (err error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	reg := fmt.Sprintf(`(%s ?= ?(\S*))`, key)
	re := regexp.MustCompile(reg)
	replacement := []byte(fmt.Sprintf("%s = %s", key, newValue))
	out := re.ReplaceAll(data, replacement)
	return ioutil.WriteFile(filename, out, 0777)
}

// GetNewToken get token from database, check it
// repeat it while we can get valid token
func (c *Core) GetNewToken() (err error) {
	for {
		token, err := c.Store.PopToken()
		if err != nil {
			return err
		}
		username, err := c.CheckToken(token)
		if err != nil {
			// get next token from db and check it
			continue
		}
		// token good
		return c.SetValidToken(token, username)
	}
}
