package main

import (
	"aatrox/im-system/client_side/client"
	"fmt"
)

func main() {
	client := client.NewClient("127.0.0.1", 8888)

	if client == nil {
		fmt.Printf("connect failed\n")
		return
	}

	fmt.Printf("connect success\n")
}
