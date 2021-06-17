package main

import (
	"fmt"
	ws "github.com/huecester/socket_server/pkg/websocket"
)

type handler struct{}

func (h handler) OnConnect(cl *ws.Client) {
	fmt.Println("Connection from", cl.IP)
	cl.Send("Hello, client!")
}

func (h handler) OnMessage(cl *ws.Client, msg string) {
	fmt.Println("Message:", msg)
	cl.Send(fmt.Sprintf("You said: %v\n", msg))
}

func (h handler) OnClose(cl *ws.Client) {
	fmt.Println("Connection closed", cl.IP)
}

func main() {
	handlerObject := handler{}
	server := ws.New(handlerObject)
	server.Start(8000)
}
