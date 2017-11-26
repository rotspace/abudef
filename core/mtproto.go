// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"log"
	"math/rand"

	"github.com/rotspace/mtproto"
)

// MTProtoManager represent mtproto client
type MTProtoManager interface {
	SendMessageToUsername(username string, message string) (err error)
	SendMessage(user UserDialog, message string) (err error)
	GetUnreadedDialogs() ([]UserDialog, error)
}

// UserDialog shows minimal params to send message
type UserDialog struct {
	ID         int32
	AccessHash int64
}

var (
	_ MTProtoManager = (*defaultMtprotoWrapper)(nil)
)

type defaultMtprotoWrapper struct {
	*mtproto.MTProto
}

// NewDefaultMTPRotoWrapper create and starts mtproto client
func NewDefaultMTPRotoWrapper() (_ MTProtoManager, err error) {

	wr := new(defaultMtprotoWrapper)

	m, err := mtproto.New(mtproto.OptAuthFile(DefaultAuthDataFilename))
	if err != nil {
		panic(err)
	}

	err = m.Connect()
	if err != nil {
		return
	}

	wr.MTProto = m

	return wr, nil
}

func (dmw *defaultMtprotoWrapper) GetUnreadedDialogs() (dialogs []UserDialog, err error) {
	resp, err := dmw.MTProto.MessagesGetDialogs()
	if err != nil {
		return
	}
	dialogs = getUnreaded(resp)
	return
}

func (dmw *defaultMtprotoWrapper) SendMessage(ud UserDialog, msg string) (err error) {
	_, err = dmw.MessagesSendMessage(false, false, false, true, mtproto.TL_inputPeerUser{
		User_id:     ud.ID,
		Access_hash: ud.AccessHash,
	}, 0, msg, rand.Int63(), mtproto.TL_null{}, nil)
	return
}

func getMessage(messages []mtproto.TL, id int32) (mtproto.TL_message, bool) {
	for _, msgTL := range messages {
		switch msg := msgTL.(type) {
		case mtproto.TL_message:
			if msg.Id == id {
				return msg, true
			}
		}
	}
	return mtproto.TL_message{}, false
}

func getUser(users []mtproto.TL, id int32) (mtproto.TL_user, bool) {
	for _, userTL := range users {
		switch user := userTL.(type) {
		case mtproto.TL_user:
			if user.Id == id {
				return user, true
			}
		}
	}
	return mtproto.TL_user{}, false
}

func getUnreaded(resp mtproto.TL_messages_dialogsSlice) (res []UserDialog) {
	for _, diaTL := range resp.Dialogs {
		dia := diaTL.(mtproto.TL_dialog)
		switch peer := dia.Peer.(type) {
		case mtproto.TL_peerUser:
			_ = peer
			if dia.Unread_count > 0 {
				user, has := getUser(resp.Users, peer.User_id)
				if !has {
					log.Println("user not found")
					continue
				}
				if user.Bot {
					log.Println(user.Username, "is bot, continue")
					continue
				}
				//log.Printf("Dialog (%d) %d %d %s", peer.User_id, dia.Unread_count, dia.Top_message, reflect.TypeOf(dia.Peer))
				_, has = getMessage(resp.Messages, dia.Top_message)
				if !has {
					log.Println("message not found")
					continue
				}
				//log.Println(msg.Message)
				res = append(res, UserDialog{ID: peer.User_id, AccessHash: user.Access_hash})
			}
		}
	}
	return
}
