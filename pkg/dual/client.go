package dual

import (
	"github.com/huecester/socket_server/pkg/tcp"
	"github.com/huecester/socket_server/pkg/ws"
)

// Client
type Client struct {
	isTCP bool
	ID string

	tcpClient *tcp.Client
	wsClient *ws.Client
}

// Constructors
func newTCP(cl *tcp.Client) Client {
	return Client{
		isTCP: true,
		ID: cl.ID,
		tcpClient: cl,
	}
}

func newWS(cl *ws.Client) Client {
	return Client{
		isTCP: false,
		ID: cl.ID,
		wsClient: cl,
	}
}

// Methods
func (cl *Client) receive() {
	if cl.isTCP {
		cl.tcpClient.Receive()
	} else {
		cl.wsClient.Receive()
	}
}

func (cl *Client) Send(msg string) {
	if cl.isTCP {
		cl.tcpClient.Send(msg)
	} else {
		cl.wsClient.Send(msg)
	}
}

func (cl *Client) Close() {
	if cl.isTCP {
		cl.tcpClient.Close()
	} else {
		cl.wsClient.Close()
	}
}
