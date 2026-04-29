package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, topic string, groupID string) *Consumer {
	log.Println(brokers[0])
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})

	return &Consumer{reader: r}
}

func (c *Consumer) ReadMessages(ctx context.Context, handler func([]byte)) {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("kafka read error: %v", err)
			continue
		}

		handler(msg.Value)
	}
}