package main

import (
	"fmt"

	"github.com/huecester/socket_server/pkg/tcp"
)

type handler struct {}

func (h handler) OnConnect(cl *tcp.Client) {
	fmt.Printf("New connection %v.\n", cl.ID)
}

func (h handler) OnMessage(cl *tcp.Client, msg string) {
	fmt.Printf("%v: %v\n", cl.ID, msg)
	cl.Send(msg)
}

func (h handler) OnClose(id string) {
	fmt.Printf("Closed connection %v.\n", id)
}


func main() {
	h := handler{}
	s := tcp.NewServer(h)
	s.Start(42069)
}
