package client

import (
	"fmt"
	"net"
)

type Client struct {
	ServerIP   string
	ServerPort int64
	Name       string
	Conn       net.Conn
}

func NewClient(serverIP string, serverPort int64) (client *Client) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIP, serverPort))

	if err != nil {
		fmt.Printf("net.Dial err: %v\n", err)
		return nil
	}

	client = &Client{
		ServerIP:   serverIP,
		ServerPort: serverPort,
		Conn:       conn,
	}
	return
}
