package logic

import (
	"context"
	"errors"
	"sync"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

var (
	uid    int
	mutex  sync.Mutex
	System = &User{}
)

type User struct {
	UID            int           `json:"uid"`
	Nickname       string        `json:"nickname"`
	EnterAt        time.Time     `json:"enter_at"`
	Addr           string        `json:"addr"`
	MessageChannel chan *Message `json:"-"`

	conn *websocket.Conn
}

func NewUser(conn *websocket.Conn, nickname string, addr string) *User {
	return &User{
		UID:            getuid(),
		Nickname:       nickname,
		EnterAt:        time.Now(),
		Addr:           addr,
		conn:           conn,
		MessageChannel: make(chan *Message),
	}
}
func getuid() int {
	mutex.Lock()
	uid++
	mutex.Unlock()
	return uid
}
func (user *User) SendMessage(ctx context.Context) {
	for msg := range user.MessageChannel {
		wsjson.Write(ctx, user.conn, msg)
	}
}

func (u *User) ReceiveMessage(ctx context.Context) error {
	var (
		receiveMsg map[string]string
		err        error
	)

	for {
		err = wsjson.Read(ctx, u.conn, &receiveMsg)
		if err != nil {
			// 判断连接是否正常关闭
			var closeErr websocket.CloseError
			if errors.As(err, &closeErr) {
				return nil
			}
			return err
		}
		// 内容发送到聊天室
		sendMsg := NewMessage(u, receiveMsg["content"])
		Broadcaster.Broadcast(sendMsg)
	}
}
