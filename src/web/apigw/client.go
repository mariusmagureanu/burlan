package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mariusmagureanu/burlan/src/pkg/log"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 2 << 8

	messageTopic = "message-trip"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	brokers []string
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	messages chan []byte

	mq *MQ

	uid string
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readFromWebSocket(ctx context.Context) {
	defer func() {
		c.hub.unregister <- c
		_ = c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { _ = c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Warning(err)
			}
			break
		}

		err = c.mq.writeToKafka(ctx, c.uid, message)

		if err != nil {
			log.Error("could not write to kafka:", err.Error())
		}

		//dest.messages <- message
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writeToWebSocket() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.messages:
			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Error("could not write to the websocket," + err.Error())
				return
			}
			_, _ = w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.messages)
			for i := 0; i < n; i++ {
				_, _ = w.Write(newline)
				_, _ = w.Write(<-c.messages)
			}

			if err := w.Close(); err != nil {
				log.Error("could not close websocket," + err.Error())
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) close() error {
	close(c.messages)
	err := c.mq.Close()

	if err != nil {
		return err
	}

	_ = c.conn.Close()
	return nil
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	u, _, ok := r.BasicAuth()
	if !ok {
		log.Warning("error parsing basic auth, attempt on query params")
	}

	if u == "" {
		u = r.URL.Query().Get("user")

		if u == "" {
			log.Warning("couldn't fetch a user, will stop registration now")
			return
		}
	}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("error when upgrading connection to ws:" + err.Error())
		return
	}

	client := &Client{hub: hub, conn: conn, messages: make(chan []byte, 2<<8)}

	mq := MQ{}
	mq.Init(hub, u, brokers)

	client.mq = &mq
	client.uid = u
	client.hub.register <- client

	go client.writeToWebSocket()
	go client.readFromWebSocket(context.Background())
	go client.mq.readFromKafka(context.Background(), u)
}
