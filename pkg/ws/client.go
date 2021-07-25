package ws

import (
	"fmt"
	"net"
	"io"

	"log"
)

const (
	// fin, rsv, mask, payload len, mask key
	minFrameLen = 6;
)

// Client
type Client struct {
	conn *net.Conn
	msgChan chan<- string
	ID string

	// Websocket specific
	continued string
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

		// Close connection if frame is invalid (below min length)
		if len(hb) < minFrameLen {
			return
		}

		f := newFrame(hb)

		switch f.opcode {
		case 0x0:
			// Continue
			if !f.fin {
				cl.continued += f.decode()
			} else {
				cl.msgChan <- fmt.Sprintf("%v:%v", cl.ID, cl.continued + f.decode())
				cl.continued = ""
			}

		case 0x1:
			// UTF-8
			if f.fin {
				cl.msgChan <- fmt.Sprintf("%v:%v", cl.ID, f.decode())
			} else {
				cl.continued = f.decode()
			}

		case 0x2:
			// Binary TODO
			log.Println("Binary not implemented")
			return

		case 0x8:
			// Close
			return

		case 0x9:
			// Ping
			ping := frame{
				fin: true,
				rsv: []bool{false, false, false},
				opcode: 0xa,

				payloadLen: f.payloadLen,
				payload: f.payload,
			}
			conn.Write(ping.encode())

		/*case 0xa:
			// Pong TODO
		*/
		}
	}
}

func (cl *Client) Send(msg string) {
	conn := *cl.conn

	f := frame{
		fin: true,
		rsv: []bool{false, false, false},
		opcode: 0x1,

		payloadLen: uint64(len(msg)),
		payload: []byte(msg),
	}

	conn.Write(f.encode())
}

func (cl *Client) Close() {
	conn := *cl.conn

	f := frame{
		fin: true,
		rsv: []bool{false, false, false},
		opcode: 0x8,
	}

	conn.Write(f.encode())
	conn.Close()

	cl.msgChan <- fmt.Sprintf("CLOSE:%v", cl.ID)
}
