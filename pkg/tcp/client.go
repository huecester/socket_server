package tcp

import (
	"fmt"
	"net"
	"io"

	"log"
)

// Client
type Client struct {
	conn *net.Conn
	msgChan chan<- string

	ID string
}

// Constructor
func newClient(conn *net.Conn, msgChan chan<- string, id string) Client {
	return Client{
		conn: conn,
		msgChan: msgChan,
		ID: id,
	}
}

// Methods
func (cl *Client) Receive() {
	buf := make([]byte, 2048)
	conn := *cl.conn

	defer cl.Close()

	for {
		n, err := conn.Read(buf)

		if err != nil {
			if err == io.EOF {
				break
			}

			log.Println(err)
			break
		}

		hb := make([]byte, n)
		copy(hb, buf)
		cl.msgChan <- fmt.Sprintf("%v:%v", cl.ID, string(hb))
	}
}

func (cl *Client) Send(msg string) {
	conn := *cl.conn
	conn.Write([]byte(msg))
}

func (cl *Client) Close() {
	conn := *cl.conn
	conn.Close()
	cl.msgChan <- fmt.Sprintf("CLOSE:%v", cl.ID)
}
