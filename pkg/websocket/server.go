package websocket

import (
	"fmt"
	"net/http"
)

// Interfaces
type handler interface {
	OnConnect(*Client)
	OnMessage(*Client, string)
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
		clients: make(map[string]*Client),
	}
}

// Methods
func (s *server) Start(port int) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		cl, err := s.newClient(w, r)
		if err != nil {
			fmt.Println(err)
		}

		s.clients[cl.IP] = &cl
		go s.handler.OnConnect(&cl)
		go cl.handleFrame()
	})

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		fmt.Println(err)
		return
	}
}
