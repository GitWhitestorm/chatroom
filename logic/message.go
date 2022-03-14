package logic

import "time"

type Message struct {
	User    *User     `json:"user"`
	Type    int       `json:"type"`
	Content string    `json:"content"`
	MsgTime time.Time `json:"msg_time"`
	Users   map[string]*User
}

const (
	MsgTypeNormal   = iota //普通用户信息
	MsgTypeSystem          //系统消息
	MsgTypeError           //错误信息
	MsgTypeUserList        //发送当前用户列表
)

func NewErrorMessage(msg string) *Message {
	return &Message{
		User:    System,
		Type:    MsgTypeError,
		Content: msg,
		MsgTime: time.Now(),
	}
}

func NewMessage(user *User, content string) *Message {
	return &Message{
		User:    user,
		Type:    MsgTypeNormal,
		Content: content,
		MsgTime: time.Now(),
	}
}
func NewWelcomeMessage(content string) *Message {
	return &Message{
		Type:    MsgTypeSystem,
		User:    System,
		Content: content,
		MsgTime: time.Now(),
	}
}

func NewNoticeMessage(content string) *Message {
	return &Message{
		Type:    MsgTypeSystem,
		User:    System,
		Content: content,
		MsgTime: time.Now(),
	}
}
