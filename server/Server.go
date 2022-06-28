package server

import (
	"aatrox/im-system/user"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
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

	server.UserOnline(user)

	//accept client input
	go server.HandleUserInput(conn, user)
}

func (server *Server) Broadcast(user *user.User, msg string) {
	sendMsg := fmt.Sprintf("[%s]%s: %s", time.Now().Format("2006-01-02 15:04:05"), user.Name, msg)

	server.Msg <- sendMsg
}

func (server *Server) SendMessage(user *user.User, msg string) {
	user.Chan <- msg
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

func (server *Server) HandleUserInput(conn net.Conn, user *user.User) {
	input := make([]byte, 4096)

	for {
		n, err := conn.Read(input)

		if n == 0 {
			server.UserOffline(user)
			return
		}

		if err != nil && err != io.EOF {
			fmt.Printf("conn.Read err: %v\n", err)
			return
		}

		// clean "\n" after user input string
		msg := string(input[:n-1])

		server.HandleMessage(user, msg)
	}
}

func (server *Server) UserOnline(user *user.User) {
	// add user to server.OnlineUsers
	server.Lock.Lock()
	server.OnlineUsers[user.Name] = user
	server.Lock.Unlock()

	// broadcast online message
	server.Broadcast(user, "online")
}

func (server *Server) UserOffline(user *user.User) {
	// add user to server.OnlineUsers
	server.Lock.Lock()
	delete(server.OnlineUsers, user.Name)
	server.Lock.Unlock()

	// broadcast offline message
	server.Broadcast(user, "offline")
}

func (server *Server) HandleMessage(user *user.User, msg string) {
	switch msg {
	case "who":
		for _, u := range server.OnlineUsers {
			output := fmt.Sprintf("%s is online", u.Name)
			server.SendMessage(user, output)
		}
		break
	default:
		server.Broadcast(user, msg)
	}
}
