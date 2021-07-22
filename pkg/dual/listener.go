package dual

import (
	"github.com/huecester/socket_server/pkg/tcp"
	"github.com/huecester/socket_server/pkg/ws"
)

// Listener
type listener struct {
	tcpListener tcp.Listener
	tcpChan <-chan *tcp.Client

	wsListener ws.Listener
	wsChan <-chan *ws.Client

	clientChan chan<- *Client
}

// Constructor
func newListener(clientChan chan<- *Client, msgChan chan<- string) listener {
	tcpChan := make(chan *tcp.Client)
	wsChan := make(chan *ws.Client)

	var readTCPChan <-chan *tcp.Client = tcpChan
	var readWSChan <-chan *ws.Client = wsChan

	var writeTCPChan chan<- *tcp.Client = tcpChan
	var writeWSChan chan<- *ws.Client = wsChan

	return listener{
		tcpListener: tcp.NewListener(writeTCPChan, msgChan),
		tcpChan: readTCPChan,

		wsListener: ws.NewListener(writeWSChan, msgChan),
		wsChan: readWSChan,

		clientChan: clientChan,
	}
}

// Methods
func (l *listener) start(tcpPort, wsPort int) {
	go l.tcpListener.Start(tcpPort)
	go l.wsListener.Start(wsPort)

	for {
		select {
		case tcpClient := <-l.tcpChan:
			nc := newTCP(tcpClient)
			l.clientChan <- &nc

		case wsClient := <-l.wsChan:
			nc := newWS(wsClient)
			l.clientChan <- &nc
		}
	}
}
