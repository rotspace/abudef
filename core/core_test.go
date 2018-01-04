// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import "testing"

func TestGetFileValue(t *testing.T) {
	var data = []byte(`APP_ENV=local
APP_KEY=base64:pnxqHgUu4Hn0FvbvLZO7JzDBfsQC3ghq6uwU5v7s+Ns=
APP_DEBUG=true
APP_LOG_LEVEL=debug

DB_CONNECTION=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=admin_bot
DB_USERNAME=admin_bot
DB_PASSWORD=tXDpm9vNNP

BROADCAST_DRIVER=log
CACHE_DRIVER=file
SESSION_DRIVER=file
QUEUE_DRIVER=sync

REDIS_HOST=127.0.0.1
REDIS_PASSWORD=null
REDIS_PORT=6379

MAIL_DRIVER=smtp
MAIL_HOST=smtp.yandex.ru
MAIL_PORT=587
MAIL_USERNAME=info@mail.ru
MAIL_PASSWORD=123456
MAIL_FROM=info@mail.ru
MAIL_NAME=NAME

PUSHER_KEY=
PUSHER_SECRET=
PUSHER_APP_ID=

TELEGRAM_BOT_TOKEN = 444499532:AAGL7GnXNmp4ouUOUMLgXYrwKIo2hb3pgaw`)
	value, err := getValue(data, "TELEGRAM_BOT_TOKEN")
	if err != nil {
		t.Fatal(err)
	}
	if value != "444499532:AAGL7GnXNmp4ouUOUMLgXYrwKIo2hb3pgaw" {
		t.Fatalf(value)
	}
}
