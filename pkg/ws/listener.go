package ws

import (
	"fmt"
	"net/http"
	"strings"

	"io"
	"crypto/sha1"
	"encoding/base64"

	"github.com/huecester/socket_server/pkg/id"

	"log"
)

// Helper functions
func getAcceptHash(key string) string {
	const uid = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	h := sha1.New()

	io.WriteString(h, key)
	io.WriteString(h, uid)
	hash := h.Sum(nil)
	accept := base64.StdEncoding.EncodeToString(hash)
	return accept
}


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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Hijack connection
		hj, ok := w.(http.Hijacker)
		if !ok {
			log.Println("Hijacking is not supported.")
			return
		}

		conn, _, err := hj.Hijack()
		if err != nil {
			log.Println(err)
			return
		}

		// Create response
		res := []string{
			"HTTP/1.1 101 Switching Protocols",
			"Upgrade: websocket",
			"Connection: upgrade",
		}

		res = append(res, "Sec-Websocket-Accept: " + getAcceptHash(r.Header.Get("Sec-Websocket-Key")))

		// Send response
		res = append(res, "", "")
		resStr := strings.Join(res, "\r\n")
		conn.Write([]byte(resStr))

		nc := newClient(&conn, l.msgChan, id.New(8))
		l.clientChan <- &nc
	})

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Println(err)
	}
}
