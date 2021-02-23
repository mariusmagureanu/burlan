package main

import (
	"context"
	"github.com/mariusmagureanu/burlan/src/pkg/auth"
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
	maxMessageSize = 2 << 12

	messageTopic = "message-trip"
)

var (
	newline = []byte{'\n'}

	upgrader = websocket.Upgrader{
		ReadBufferSize:  2 << 11,
		WriteBufferSize: 2 << 11,
	}

	brokers []string
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub      *Hub
	mq       *MQ
	conn     *websocket.Conn
	messages chan []byte
	uid      string
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readFromWebSocket(ctx context.Context, claim *auth.JwtClaim) {
	defer func() {
		c.hub.unregister <- c
		_ = c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { _ = c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		err := claim.Valid()
		if err != nil {
			log.Warning("claim error while reading ws, ", err.Error())
			break
		}

		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Warning(err)
			} else {
				log.Error(err)
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
			log.Debug("writing message to websocket:", string(message))
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

	authToken := r.Header.Get("X-JWT")

	if authToken == "" {
		authToken = r.URL.Query().Get("jwt")
	}

	if authToken == "" {
		log.Warning("jwt not provided, request will not continue")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	claim, err := jwtWrapper.ValidateToken(authToken)

	if err != nil {
		log.Warning("invalid jwt ", err.Error())
		w.WriteHeader(http.StatusForbidden)
		return
	}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("error when upgrading connection to ws:" + err.Error())
		return
	}

	client := &Client{hub: hub, conn: conn, messages: make(chan []byte, 2<<8)}

	mq := MQ{}
	mq.Init(hub, claim.ClientUID, brokers)

	client.mq = &mq
	client.uid = claim.ClientUID
	client.hub.register <- client

	go client.mq.readFromKafka(context.Background(), claim.ClientUID)
	go client.writeToWebSocket()
	go client.readFromWebSocket(context.Background(), claim)
}
