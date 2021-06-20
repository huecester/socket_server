package main

import (
	"fmt"

	"github.com/huecester/socket_server/pkg/tcp"
	//ws "github.com/huecester/socket_server/pkg/websocket"
)

const (
	tcpPort = 42069
)

type handlerImpl struct{}

func (h handlerImpl) OnConnect(cl *tcp.Client) {
	fmt.Println("Conneciton opened.")
	cl.Send("Hello, client!")
}

func (h handlerImpl) OnMessage(cl *tcp.Client, msg string) {
	fmt.Println("Received:", msg)
}

func (h handlerImpl) OnClose(cl *tcp.Client) {
	fmt.Println("Connection closed.")
}


func main() {
	var h handlerImpl
	server := tcp.New(h)

	err := server.Start(tcpPort)
	if err != nil {
		fmt.Println(err)
		return
	}
}
