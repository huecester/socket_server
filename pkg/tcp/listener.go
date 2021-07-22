package tcp

import (
	"fmt"
	"net"

	"github.com/huecester/socket_server/pkg/id"

	"log"
)

// Listener
type Listener struct {
	clientChan chan<- *Client
	msgChan chan<- string
}

// Constructor
func NewListener(clientChan chan<- *Client, msgChan chan<- string) Listener {
	return Listener{
		clientChan: clientChan,
		msgChan: msgChan,
	}
}

// Methods
func (l *Listener) Start(port int) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Println(err)
		return
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		nc := newClient(&conn, l.msgChan, id.New(8))
		l.clientChan <- &nc
	}
}
