package logic

import "log"

type broadcaster struct {
	users                 map[string]*User
	enteringChannel       chan *User
	leavingChannel        chan *User
	messageChannel        chan *Message
	checkUserChannel      chan string
	checkUserCanInChannel chan bool
}

var Broadcaster = &broadcaster{
	users:                 make(map[string]*User),
	enteringChannel:       make(chan *User),
	leavingChannel:        make(chan *User),
	messageChannel:        make(chan *Message),
	checkUserChannel:      make(chan string),
	checkUserCanInChannel: make(chan bool),
}

func (b *broadcaster) CanEnterRoom(nickname string) bool {
	log.Println(nickname)
	b.checkUserChannel <- nickname

	return <-b.checkUserCanInChannel
}

func (b *broadcaster) Broadcast(msg *Message) {
	b.messageChannel <- msg
}
func (b *broadcaster) Start() {
	for {
		select {
		case msg := <-b.messageChannel:
			for _, user := range b.users {
				if user.UID == msg.User.UID {
					continue
				}
				user.MessageChannel <- msg
			}
		case user := <-b.leavingChannel:
			delete(b.users, user.Nickname)
		case user := <-b.enteringChannel:
			b.users[user.Nickname] = user
		case nickname := <-b.checkUserChannel:
			if _, ok := b.users[nickname]; !ok {
				b.checkUserCanInChannel <- true
			} else {
				b.checkUserCanInChannel <- false
			}

		}
	}
}

func (b *broadcaster) UserEntering(user *User) {
	b.enteringChannel <- user
}

func (b *broadcaster) UserLeaving(user *User) {
	b.leavingChannel <- user
}
