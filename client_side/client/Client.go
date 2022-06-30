package client

import (
	"fmt"
	"io"
	"net"
	"os"
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
	go client.ReadMsgFromServer()

	for client.flag != 0 {
		for client.ShowMenu() != true {
		}

		switch client.flag {
		case 1:
			client.PublicChat()
		case 2:
			client.PrivateChat()
		case 3:
			client.Rename()
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

func (client *Client) ReadMsgFromServer() {
	io.Copy(os.Stdout, client.Conn)

	// same as

	// for {
	// 	msg := make([]byte, 4096)
	// 	client.Conn.Read(msg)

	// 	fmt.Print(string(msg))
	// }
}

func (client *Client) SendMsgToServer(msg string) (ok bool) {
	ok = true

	_, err := client.Conn.Write([]byte(fmt.Sprintf("%s\n", msg)))
	if err != nil {
		fmt.Printf("client.Conn.Write err: %v", err)
		ok = false
	}
	return
}

func (client *Client) PublicChat() (ok bool) {
	var inputMsg string
	fmt.Println("Insert message, \"exit\" to leave:")
	fmt.Scanln(&inputMsg)

	for inputMsg != "exit" {

		if len(inputMsg) != 0 {
			ok = client.SendMsgToServer(inputMsg)
			if !ok {
				break
			}
		}

		inputMsg = ""
		fmt.Println("Insert message, \"exit\" to leave:")
		fmt.Scanln(&inputMsg)
	}
	return
}

func (client *Client) PrivateChat() (ok bool) {
	var targetName string
	var msg string

	ok = client.ListUsers()
	fmt.Println("Insert target username, \"exit\" to leave:")
	fmt.Scanln(&targetName)

	if targetName != "exit" {
		fmt.Println("Insert message, \"exit\" to leave:")
		fmt.Scanln(&msg)

		for msg != "exit" {
			msg = fmt.Sprintf("dm:%s %s", targetName, msg)
			client.SendMsgToServer(msg)

			msg = ""
			fmt.Println("Insert message, \"exit\" to leave:")
			fmt.Scanln(&msg)
		}
	}

	return
}

func (client *Client) ListUsers() (ok bool) {
	ok = client.SendMsgToServer("who")
	return
}

func (client *Client) Rename() (ok bool) {
	fmt.Print("Input new name: ")
	fmt.Scanln(&client.Name)

	msg := fmt.Sprintf("rename %s", client.Name)
	ok = client.SendMsgToServer(msg)

	return
}
