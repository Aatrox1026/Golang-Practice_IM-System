package main

import "aatrox/im-system/server_side/server"

func main() {
	server := server.NewServer("127.0.0.1", 8888)
	server.Start()
}
