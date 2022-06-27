package server

import (
	"aatrox/im-system/user"
	"fmt"
	"net"
	"sync"
)

type Server struct {
	IP   string
	Port int64

	OnlineUsers map[string]*user.User
	Lock        sync.RWMutex

	Msg chan string
}

// create server
func NewServer(ip string, port int64) (server *Server) {
	server = &Server{
		IP:          ip,
		Port:        port,
		OnlineUsers: make(map[string]*user.User),
		Msg:         make(chan string),
	}

	return
}

// start server
func (server *Server) Start() {
	// socket listen
	address := fmt.Sprintf("%s:%d", server.IP, server.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("net.Listen err: %v\n", err)
		return
	}

	// close socket listen
	defer listener.Close()

	// listen message
	go server.ListenMessage()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("listener.Accept err: %v\n", err)
			continue
		}

		//do handler
		go server.Handle(conn)
	}

}

func (server *Server) Handle(conn net.Conn) {
	// fmt.Println("create conn success")

	// user online
	user := user.NewUser(conn)

	// add user to server.OnlineUsers
	server.Lock.Lock()
	server.OnlineUsers[user.Name] = user
	server.Lock.Unlock()

	// broadcast online message
	server.Broadcast(user, "go online")
}

func (server *Server) Broadcast(user *user.User, msg string) {
	sendMsg := fmt.Sprintf("[%s]%s: %s", user.Addr, user.Name, msg)

	server.Msg <- sendMsg
}

func (server *Server) ListenMessage() {
	for {
		msg := <-server.Msg

		server.Lock.Lock()
		for _, user := range server.OnlineUsers {
			user.Chan <- msg
		}
		server.Lock.Unlock()
	}
}
