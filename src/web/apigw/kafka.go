package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mariusmagureanu/burlan/src/pkg/log"
	"github.com/segmentio/kafka-go"
	"io"
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
			log.Warning("could not read kafka message,", err.Error())
			continue
		}

		// this check is probably redundant, as each consumer is alone
		// in its own group, so it's (I guess) guaranteed that each message
		// reaches the correct destination.
		destClientKey := string(msg.Key)
		log.Debug("got message,", string(msg.Value))
		if destClientKey != toUID {
			log.Error("waat??", destClientKey, toUID)
			continue
		}

		if dest, ok := mq.hub.clients[destClientKey]; ok {
			var im internalMessage
			err = json.Unmarshal(msg.Value, &im)

			if err != nil {
				log.Warning("received invalid message format")
				continue
			}

			dest.messages <- msg.Value
			err = mq.kafkaReader.CommitMessages(ctx, msg)

			if err != nil {
				log.Error("error when committing message to kafka,", err.Error())
			}
		}
	}
}

func (mq *MQ) writeToKafka(ctx context.Context, fromUID string, message []byte) error {

	//message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

	msg := kafka.Message{}
	fmt.Println(fmt.Sprintf("%s", message))

	// we'll need this later for broadcasting to groups
	//c.hub.broadcast <- message
	var im internalMessage

	//tmpMsg := fmt.Sprintf("%s", message)

	err := json.Unmarshal(message, &im)
	if err != nil {
		log.Error("unmarshalling error:", err.Error())
		return err
	}

	im.From = fromUID

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

	log.Info(fmt.Sprintf("sent message <%s> from <%s> to <%s>", im.Text, im.From, im.To))
	return nil
}
