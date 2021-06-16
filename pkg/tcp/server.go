package tcp

import (
	"fmt"
	"io"
	"net"

	idgen "github.com/huecester/socket_server/pkg/id"
)

// Interfaces
type handler interface {
	OnConnect(*Client)
	OnMessage(*Client, string)
}

////////////
// Client //
////////////

type Client struct {
	server *server
	conn   *net.Conn
	id     string
}

// Constructor
func newClient(serverPtr *server, connPtr *net.Conn, id string) Client {
	client := Client{
		server: serverPtr,
		conn:   connPtr,
		id:     id,
	}

	return client
}

// Methods
func (cl *Client) receive() {
	buf := make([]byte, 1024)
	conn := cl.conn

	for {
		n, err := conn.Read(buf)

		if err != nil {
			if err == io.EOF {
				continue
			}

			fmt.Println(err)

			break
		}

		hb := make([]byte, n)
		copy(hb, buf)
		go cl.s.handler.OnMessage(cl, string(hb))
	}
}

func (cl *Client) Send(message string) {
	conn := *cl.conn
	conn.Write([]byte(message))
}

func (cl *Client) Close() {
	conn := *cl.conn
	conn.Close()
	delete(cl.server.clients, cl.id)
}

////////////
// Server //
////////////

type server struct {
	handler handler
	clients map[string]*Client
}

// Constructor
func New(handler handler) server {
	return server{
		handler: handler,
	}
}

// Methods
func (s *server) Start(port int) error {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	defer ln.Close()

	for {
		select {}
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		id := idgen.New(8)
		cl := newClient(s, &conn, id)

		s.clients[id] = &cl
		go s.handler.OnConnect(&cl)
		go cl.receive()
	}

	return nil
}
