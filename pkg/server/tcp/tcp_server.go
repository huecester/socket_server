package tcp

import (
	"fmt"
	"io"
	"net"
)

/////////////////
// Connection  //
/////////////////
type Client struct {
	server *server
	conn   *net.Conn
}

// Constructor
func newClient(serverPtr *server, connPtr *net.Conn) Client {
	client := Client{
		server: serverPtr,
		conn:   connPtr,
	}

	return client
}

// Methods
func (cl Client) Send(message string) {
	conn := *cl.conn
	conn.Write([]byte(message))
}

////////////
// Server //
////////////

// Helper types
type handler interface {
	OnConnect(*Client)
	OnMessage(*Client, string)
}

type server struct {
	handler handler
}

// Constructor
func New(handler handler) server {
	return server{
		handler: handler,
	}
}

// Methods
func (s server) receive(cl *Client) {
	buffer := make([]byte, 1024)
	conn := *cl.conn

	for {
		n, err := conn.Read(buffer)

		if err != nil {
			if err == io.EOF {
				continue
			}

			fmt.Println(err)

			break
		}

		hb := make([]byte, n)
		copy(hb, buffer)
		go s.handler.OnMessage(cl, string(hb))
	}
}

func (s *server) Start(port int) error {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		cl := newClient(s, &conn)

		go s.handler.OnConnect(&cl)
		go s.receive(&cl)
	}

	return nil
}
