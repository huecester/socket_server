package websocket

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"

	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"strings"
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

func getIP(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

///////////
// Frame //
///////////

type frame struct {
	raw []byte

	fin    bool
	rsv    []bool
	opcode byte
	mask   bool

	payloadLen uint64
	maskKey    []byte
	payload    []byte

	decoded string
}

// Constructor
func newFrame(data []byte) (frame, error) {
	// Save raw data
	raw := make([]byte, len(data))
	copy(raw, data)

	// Header
	// fin
	fin := data[0]&(1<<7) > 0

	// rsv bits
	rsv := make([]bool, 3)
	for i := 0; i < 3; i++ {
		rsv[i] = data[0]&byte(1<<(6-i)) > 0
	}

	// opcode
	opcode := data[0] & 0b00001111

	data = data[1:]

	// mask
	mask := data[0]&(1<<7) > 0

	// Data
	// Payload length
	var payloadLen uint64
	payloadByte := data[0] & 0b01111111
	data = data[1:]

	if payloadByte := int(payloadByte); payloadByte <= 125 {
		payloadLen = uint64(payloadByte)
	} else {
		// Extended payload
		payloadByteSlice := make([]byte, 8)
		if payloadByte < 127 {
			payloadByteSlice = data[:2]
			data = data[2:]
		} else {
			payloadByteSlice = data[:8]
			data = data[8:]
		}

		// Binary to uint
		payloadLen = binary.BigEndian.Uint64(payloadByteSlice)
	}

	// Mask key
	maskKey := make([]byte, 4)
	if mask {
		maskKey = data[:4]
		data = data[4:]
	}

	// Additional
	// Decoded payload
	decodedBytes := make([]byte, 0, payloadLen)
	var i uint64
	for i = 0; i < payloadLen; i++ {
		decodedBytes = append(decodedBytes, data[i]^maskKey[i%4])
	}
	decoded := string(decodedBytes)

	return frame{
		raw:        raw,
		fin:        fin,
		rsv:        rsv,
		opcode:     opcode,
		mask:       mask,
		payloadLen: payloadLen,
		maskKey:    maskKey,
		payload:    data,
		decoded:    decoded,
	}, nil
}

func createMessageFrame(message string) []byte {
	res := []byte{
		0b10000001,         // fin true | rsv{0, 0, 0} | opcode 0x01
		byte(len(message)), // mask false | payload length
	}

	for _, r := range message {
		res = append(res, byte(r))
	}

	return res
}

////////////
// Client //
////////////

type Client struct {
	server *server
	conn   *net.Conn
	IP     string

	continueHead frame
	continueBody []frame

	sentClose bool
}

// Constructor
func (s *server) newClient(w http.ResponseWriter, r *http.Request) (Client, error) {
	hj, ok := w.(http.Hijacker)
	if !ok {
		return Client{}, errors.New("Hijacking is not supported.")
	}

	conn, _, err := hj.Hijack()
	if err != nil {
		return Client{}, err
	}

	res := []string{
		"HTTP/1.1 101 Switching Protocols",
		"Upgrade: websocket",
		"Connection: upgrade",
	}

	res = append(res, "Sec-WebSocket-Accept: "+getAcceptHash(r.Header.Get("Sec-WebSocket-Key")))
	//res = append(res, "Sec-WebSocket-Protocol: " + "bag")

	res = append(res, "", "")
	resStr := strings.Join(res, "\r\n")
	conn.Write([]byte(resStr))

	return Client{
		server: s,
		conn:   &conn,
		IP:     getIP(r),
	}, nil
}

// Methods
func (cl *Client) sendBytes(message []byte) {
	conn := *cl.conn
	conn.Write(message)
}

func (cl *Client) Send(message string) {
	conn := *cl.conn
	conn.Write(createMessageFrame(message))
}

func (cl *Client) Close() {
	cl.sendBytes([]byte{
		0b10001000,
		0b00000000,
	})
	cl.sentClose = true
}

func (cl *Client) Terminate() {
	conn := *cl.conn
	conn.Close()
	go cl.server.handler.OnClose(cl)
}

func (cl *Client) receiveFrame() (frame, error) {
	conn := *cl.conn
	buf := make([]byte, 1024)

	_, err := conn.Read(buf)
	if err != nil {
		return frame{}, err
	}

	data, err := newFrame(buf)
	if err != nil {
		return frame{}, err
	}

	return data, nil
}

func (cl *Client) handleFrame() {
	for {
		data, err := cl.receiveFrame()
		if err != nil {
			if err == io.EOF {
				cl.Terminate()
				return
			}

			fmt.Println(err)
			continue
		}

		switch data.opcode {
		case 0x0:
			// Continue
			if data.fin {
				message := cl.continueHead.decoded
				for _, f := range cl.continueBody {
					message += f.decoded
				}

				cl.continueHead = frame{}
				cl.continueBody = make([]frame, 0)

				go cl.server.handler.OnMessage(cl, message)
			} else {
				cl.continueBody = append(cl.continueBody, data)
			}

		case 0x1:
			// Message
			if data.fin {
				go cl.server.handler.OnMessage(cl, data.decoded)
			} else {
				cl.continueHead = data
			}

		case 0x2:
			// Binary
			if data.fin {
				go cl.server.handler.OnMessage(cl, data.decoded)
			} else {
				cl.continueHead = data
			}

		case 0x8:
			// Close
			if !cl.sentClose {
				cl.Close()
			}
			conn := *cl.conn
			conn.Close()
			go cl.server.handler.OnClose(cl)

		case 0x9:
			// Ping
			pong := make([]byte, len(data.raw))
			copy(pong, data.raw)
			pong[0] = (pong[0] & 0b11110000) | 0b1010 // Clear opcode then set to 0xa (pong)
			cl.sendBytes(pong)

			//case 0xa: // TODO send pings
			// Pong

		}
	}
}
