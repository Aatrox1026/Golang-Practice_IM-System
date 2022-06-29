package main

import (
	"aatrox/im-system/client_side/client"
	"flag"
	"fmt"
)

var argIP string
var argPort int64

func init() {
	flag.StringVar(&argIP, "ip", "127.0.0.1", "set server ip")
	flag.Int64Var(&argPort, "port", 8888, "set server port")
	flag.Parse()
}

func main() {

	client := client.NewClient(argIP, argPort)

	if client == nil {
		fmt.Printf("connect failed\n")
		return
	}

	fmt.Printf("connect success\n")

	client.Run()
}
