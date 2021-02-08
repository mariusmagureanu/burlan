package main

import (
	"fmt"

	"github.com/mariusmagureanu/burlan/src/pkg/log"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[string]*Client

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]*Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.uid] = client
			log.Info(fmt.Sprintf("connected client <%s>", client.uid))
		case client := <-h.unregister:
			if _, ok := h.clients[client.uid]; ok {
				delete(h.clients, client.uid)
				err := client.close()
				if err != nil {
					log.Error(err)
					continue
				}
				log.Info(fmt.Sprintf("disconnected client <%s>", client.uid))
			}
		case message := <-h.broadcast:
			for uid := range h.clients {

				select {
				case h.clients[uid].messages <- message:
				default:
					close(h.clients[uid].messages)
					delete(h.clients, uid)
				}
			}
		}
	}
}
