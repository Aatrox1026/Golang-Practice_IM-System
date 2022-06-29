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
	flag       int64
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
		flag:       -1,
	}
	return
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.ShowMenu() != true {
		}

		switch client.flag {
		case 1:
			fmt.Printf("Join Chat")
		case 2:
			fmt.Printf("Private Message")
		case 3:
			fmt.Printf("Rename")
		}
	}
}

func (client *Client) ShowMenu() (ok bool) {
	var flag int64

	fmt.Println("1.Join Chat")
	fmt.Println("2.Private Message")
	fmt.Println("3.Rename")
	fmt.Println("0.Exit")

	fmt.Scanln(&flag)

	switch flag {
	case 0, 1, 2, 3:
		ok = true
		client.flag = flag
	default:
		ok = false
		fmt.Println("Illegal input")
	}

	return
}
