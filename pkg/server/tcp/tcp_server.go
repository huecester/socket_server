package tcp_server

import (
	"fmt"
	"net"
)

func Server(port int, out chan<- net.Conn) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		out <- conn
	}
}
