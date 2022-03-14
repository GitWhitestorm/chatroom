package main

import (
	"bufio"
	"fmt"
	"log"
	"sync"

	"net"
	"strconv"

	"time"
)

type Message struct {
	OwnerID int
	Content string
}

var (
	// 新用户到来，进行登记
	enteringChannel = make(chan *User)
	// 用户离开 进行登记
	leavingChannel = make(chan *User)

	// uid互斥锁
	mutex sync.Mutex
	uid   = 0
	// 用于广播的信息
	messageChannel = make(chan *Message, 8)
)

type User struct {
	ID             int
	Addr           string
	EnterAt        time.Time
	MessageChannel chan string
}

func main() {
	listener, err := net.Listen("tcp", ":2000")
	if err != nil {
		panic(err)
	}
	// 开启协程用于广播
	go broadcaster()

	// 监听是否有新用户请求加入
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		// 开启一个协程处理请求
		go handleConn(conn)

	}

}

// 记录用户 用于广播
func broadcaster() {
	users := make(map[*User]struct{})
	for {
		select {
		case user := <-enteringChannel:
			users[user] = struct{}{}
		case user := <-leavingChannel:
			delete(users, user)
			// 避免goroutine泄露
			close(user.MessageChannel)
		case msg := <-messageChannel:
			// 给所有在线用户发送信息
			for user := range users {

				if user.ID == msg.OwnerID {
					continue
				}
				user.MessageChannel <- msg.Content
			}
		}
	}
}

// 处理链接
func handleConn(conn net.Conn) {
	defer conn.Close()

	var userActive = make(chan struct{})
	go func() {
		d := 5 * time.Minute
		timer := time.NewTimer(d)
		for {
			select {
			case <-timer.C:
				conn.Close()
			case <-userActive:
				timer.Reset(d)
			}
		}
	}()
	// 构建用户实体
	user := &User{
		ID:             GenUserID(),
		Addr:           conn.RemoteAddr().String(),
		EnterAt:        time.Now(),
		MessageChannel: make(chan string, 8),
	}

	// 开启一个协程用户发送信息
	go sendMessage(conn, user.MessageChannel)

	// 给当前用户发送欢迎信息，给所有用户告知新用户到来
	user.MessageChannel <- "Welcome," + user.String()

	messageChannel <- &Message{OwnerID: user.ID, Content: "user:`" + strconv.Itoa(user.ID) + "` has enter"}

	// 将该记录到全局的用户列表中，避免用锁
	enteringChannel <- user

	// 循环读取用户的输入
	input := bufio.NewScanner(conn)

	for input.Scan() {
		messageChannel <- &Message{OwnerID: user.ID, Content: strconv.Itoa(user.ID) + ":" + input.Text()}

		userActive <- struct{}{}
	}
	if err := input.Err(); err != nil {
		log.Println("读取错误", err)
	}

	// 用户离开
	leavingChannel <- user
	messageChannel <- &Message{OwnerID: user.ID, Content: "user:`" + strconv.Itoa(user.ID) + "` has left"}

}

// <-chan 只能接收数据
func sendMessage(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintf(conn, msg)
	}
}

// 加锁
func GenUserID() int {
	mutex.Lock()
	uid++
	mutex.Unlock()
	return uid
}

func (user *User) String() string {

	return fmt.Sprintln("["+user.Addr+"]"+"ID:"+strconv.Itoa(user.ID)+"  time:", user.EnterAt)
}
