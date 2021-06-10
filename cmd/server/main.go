package main

import (
	"fmt"
	"net"

	tcp "github.com/huecester/socket_server/pkg/server/tcp"
	//ws "github.com/huecester/socket_server/pkg/server/websocket"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Fprintf(conn, "Talking to client.")
}

func main() {
	inbound := make(chan net.Conn)
	go tcp.Server(42069, inbound)

	fmt.Println("Server started.")

	for {
		go handleConnection(<-inbound)
	}
}
