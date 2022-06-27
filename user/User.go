package user

import (
	"fmt"
	"net"
)

type User struct {
	Name string
	Addr string
	Chan chan string
	Conn net.Conn
}

// create user
func NewUser(conn net.Conn) (user *User) {
	user = &User{
		Name: conn.RemoteAddr().String(),
		Addr: conn.RemoteAddr().String(),
		Chan: make(chan string),
		Conn: conn,
	}

	// start to listen channel
	go user.ListenMessage()

	return
}

// listen user.Chan and send message to client
func (user *User) ListenMessage() {
	for {
		msg := fmt.Sprintf("%s\n", <-user.Chan)
		user.Conn.Write([]byte(msg))
	}
}
