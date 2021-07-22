package main

import (
	"fmt"

	"github.com/huecester/socket_server/pkg/dual"
)

type handler struct {}

func (h handler) OnConnect(cl *dual.Client) {
	fmt.Printf("New connection %v.\n", cl.ID)
}

func (h handler) OnMessage(cl *dual.Client, msg string) {
	fmt.Printf("%v: %v\n", cl.ID, msg)
	cl.Send(msg)
}

func (h handler) OnClose(id string) {
	fmt.Printf("Closed connection %v.\n", id)
}


func main() {
	h := handler{}
	s := dual.NewServer(h)
	s.Start(42069, 42070)
}
