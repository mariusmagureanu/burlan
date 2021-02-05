package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
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
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type InternalMessage struct {
	From string
	To   string
	Text string
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	messages chan []byte

	kafkaReader *kafka.Reader
	kafkaWriter kafka.Writer

	uid string
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump(ctx context.Context) {
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
				log.Printf("error: %v\n", err)
			}
			break
		}

		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		msg := kafka.Message{}

		// we'll need this later for broadcasting to groups
		//c.hub.broadcast <- message

		payload := strings.Split(string(message), ",")
		if len(payload) != 2 {
			log.Println("invalid message format")
			continue
		}

		im := InternalMessage{}
		im.From = c.uid
		im.To = payload[0]
		im.Text = payload[1]

		imPayload, err := json.Marshal(im)
		msg.Value = imPayload
		msg.Key = []byte(im.To)

		err = c.kafkaWriter.WriteMessages(ctx, msg)
		log.Println(fmt.Sprintf("sent message [%s] from [%s] to [%s]", im.Text, im.From, im.To))

		if err != nil {
			log.Println("error writing to kafka:" + err.Error())
		}

		//dest.messages <- message
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
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
				log.Println("ws write error:" + err.Error())
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
				log.Println("ws close error:" + err.Error())
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

func (c *Client) listenKafka(ctx context.Context) {

	for {
		msg, err := c.kafkaReader.FetchMessage(ctx)
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println("could not read message " + err.Error())
			continue
		}

		// this check is probably redundant, as each consumer is alone
		// in its own group, so it's (I guess) guaranteed that each message
		// reaches the correct destination.
		destClientKey := string(msg.Key)
		if destClientKey != c.uid {
			continue
		}

		if dest, ok := c.hub.clients[destClientKey]; ok {
			var im InternalMessage
			err = json.Unmarshal(msg.Value, &im)

			if err != nil {
				log.Println("received invalid message format")
				continue
			}

			dest.messages <- []byte(fmt.Sprintf("received message [%s] from [%s]", im.Text, im.From))
			err = c.kafkaReader.CommitMessages(ctx, msg)

			if err != nil {
				log.Println("error when committing message:" + err.Error())
			}
		}
	}
}

func (c *Client) close() error {
	close(c.messages)
	err := c.kafkaWriter.Close()

	if err != nil {
		return err
	}

	err = c.kafkaReader.Close()
	if err != nil {
		return err
	}

	c.conn.Close()
	return nil
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error when upgrading connection to ws:" + err.Error())
		return
	}

	u, _, ok := r.BasicAuth()
	if !ok {
		log.Println("error parsing basic auth")
		w.WriteHeader(401)
		return
	}

	client := &Client{hub: hub, conn: conn, messages: make(chan []byte, 2<<8)}

	client.kafkaReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{"localhost:9092"},
		Topic:       messageTopic,
		StartOffset: kafka.LastOffset,
		GroupID:     u,
	})

	client.kafkaWriter = kafka.Writer{Addr: kafka.TCP("localhost:9092"), Topic: messageTopic}

	client.uid = u
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump(context.Background())
	go client.listenKafka(context.Background())
}
