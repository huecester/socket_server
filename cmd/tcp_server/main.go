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
	cl.Send("Hello, client!")
}

func (h handlerImpl) OnMessage(cl *tcp.Client, msg string) {
	fmt.Println("Received:", msg)
}

func main() {
	var handler handlerImpl
	server := tcp.New(handler)

	err := server.Start(tcpPort)
	if err != nil {
		fmt.Println(err)
		return
	}
}
