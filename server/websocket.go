package server

import (
	"log"
	"net/http"

	"chatroom/logic"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func WebSocketHandleFunc(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println("websocket accept error:", err)
		return
	}

	log.Println("用户已进入1")

	// 新用户进来，构建实例
	nickname := r.FormValue("nickname")

	if l := len(nickname); l < 2 || l > 20 {
		log.Println("nickname illegal:", nickname)
		wsjson.Write(r.Context(), conn, logic.NewErrorMessage("非法昵称，昵称长度：4-20"))
		conn.Close(websocket.StatusUnsupportedData, "nickname illegal")
		return
	}
	if !logic.Broadcaster.CanEnterRoom(nickname) {
		log.Println("昵称已经存在:", nickname)
		wsjson.Write(r.Context(), conn, logic.NewErrorMessage("该昵称已经存在"))
		conn.Close(websocket.StatusUnsupportedData, "nickname exists!")
		return
	}

	user := logic.NewUser(conn, nickname, r.RemoteAddr)
	log.Println("用户注册成功")
	// 开启给用户发送信息的goroutine
	go user.SendMessage(r.Context())

	// 给当前用户发送欢迎信息
	user.MessageChannel <- logic.NewWelcomeMessage(nickname)

	// 给所有用户告知新用户到来
	msg := logic.NewNoticeMessage(nickname + "加入了聊天室")

	logic.Broadcaster.Broadcast(msg)

	// 将该用户加入广播器的用户列表中
	logic.Broadcaster.UserEntering(user)

	log.Println("user:", nickname, "joins chat")
	// 接收用户消息
	err = user.ReceiveMessage(r.Context())

	// 用户离开
	logic.Broadcaster.UserLeaving(user)
	msg = logic.NewNoticeMessage(user.Nickname + "离开了聊天室")
	logic.Broadcaster.Broadcast(msg)
	log.Println("user:", nickname, " leaves chat")

	if err == nil {
		conn.Close(websocket.StatusNormalClosure, "")
	} else {
		log.Println("read from client error", err)
		conn.Close(websocket.StatusInternalError, "Read from client error")
	}

}
