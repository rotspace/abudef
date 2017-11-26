// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
)

var (
	// DefaultAuthDataFilename file name used for mtproto auth data
	DefaultAuthDataFilename = "authdata"
	// DefaultTokensFilename file used for tokens data
	DefaultTokensFilename = "tokens.json"
)

// Store interface
type Store interface {
	PopToken() (string, error)
	PushToken(string) error

	// SetAuthData(data []byte) error
	// GetAuthData() ([]byte, error)
}

// NewStore returns defaultStore
func NewStore() (_ Store, err error) {
	s := new(defaultStore)
	err = s.loadpersistentTokens()
	if err != nil {
		return
	}
	return s, nil
}

var _ Store = (*defaultStore)(nil)

type defaultStore struct {
	tokens   []string
	tokensMu sync.Mutex
}

func (ds *defaultStore) persistTokens() (err error) {
	data, err := json.Marshal(ds.tokens)
	if err != nil {
		return
	}
	return ioutil.WriteFile(DefaultTokensFilename, data, 0777)
}

func (ds *defaultStore) loadpersistentTokens() (err error) {
	// skip error handling
	data, _ := ioutil.ReadFile(DefaultTokensFilename)
	return json.Unmarshal(data, &ds.tokens)
}

func (ds *defaultStore) PopToken() (token string, err error) {
	ds.tokensMu.Lock()
	defer ds.tokensMu.Unlock()

	if len(ds.tokens) == 0 {
		return "", fmt.Errorf("Not enogh tokens")
	}

	token = ds.tokens[0]

	ds.tokens = ds.tokens[1:]

	return token, ds.persistTokens()
}

func (ds *defaultStore) PushToken(token string) (err error) {
	ds.tokens = append(ds.tokens, token)
	return ds.persistTokens()
}

//
// func (ds *defaultStore) SetAuthData(data []byte) (err error) {
// 	err = ioutil.WriteFile(DefaultAuthDataFilename, data, 0777)
// 	return
// }
//
// func (ds *defaultStore) GetAuthData() (data []byte, err error) {
// 	return ioutil.ReadFile(DefaultAuthDataFilename)
// }
