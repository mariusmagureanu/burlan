package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/mariusmagureanu/burlan/src/pkg/errors"
	"github.com/segmentio/kafka-go"
	"io"
	"log"
	"strings"
)

type internalMessage struct {
	From string
	To   string
	Text string
}

type MQ struct {
	kafkaReader *kafka.Reader
	kafkaWriter kafka.Writer
	hub         *Hub
}

func (mq *MQ) Init(h *Hub, groupID string, brokers []string) {
	mq.kafkaReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		Topic:       messageTopic,
		StartOffset: kafka.LastOffset,
		GroupID:     groupID,
	})

	mq.kafkaWriter = kafka.Writer{Addr: kafka.TCP(brokers...), Topic: messageTopic}

	mq.hub = h
}

func (mq *MQ) Close() error {
	err := mq.kafkaReader.Close()

	if err != nil {
		return err
	}

	return mq.kafkaWriter.Close()
}

func (mq *MQ) readFromKafka(ctx context.Context, toUID string) {

	for {
		msg, err := mq.kafkaReader.FetchMessage(ctx)
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
		if destClientKey != toUID {
			continue
		}

		if dest, ok := mq.hub.clients[destClientKey]; ok {
			var im internalMessage
			err = json.Unmarshal(msg.Value, &im)

			if err != nil {
				log.Println("received invalid message format")
				continue
			}

			dest.messages <- []byte(fmt.Sprintf("received message [%s] from [%s]", im.Text, im.From))
			err = mq.kafkaReader.CommitMessages(ctx, msg)

			if err != nil {
				log.Println("error when committing message:" + err.Error())
			}
		}
	}
}

func (mq *MQ) writeToKafka(ctx context.Context, fromUID string, message []byte) error {

	message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
	msg := kafka.Message{}

	// we'll need this later for broadcasting to groups
	//c.hub.broadcast <- message

	payload := strings.Split(string(message), ",")
	if len(payload) != 2 {
		return errors.ErrWSInvalidFormat
	}

	im := internalMessage{}
	im.From = fromUID
	im.To = payload[0]
	im.Text = payload[1]

	imPayload, err := json.Marshal(im)
	if err != nil {
		return err
	}

	msg.Value = imPayload
	msg.Key = []byte(im.To)

	err = mq.kafkaWriter.WriteMessages(ctx, msg)

	if err != nil {
		return err
	}

	log.Println(fmt.Sprintf("sent message [%s] from [%s] to [%s]", im.Text, im.From, im.To))
	return nil
}
