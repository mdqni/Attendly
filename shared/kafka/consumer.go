package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
)

type MessageHandler func(key, value []byte) error

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, topic, groupID string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			GroupID:  groupID,
			Topic:    topic,
			MinBytes: 1,
			MaxBytes: 10e6,
		}),
	}
}

func (c *Consumer) Start(ctx context.Context, handler MessageHandler) error {
	defer c.reader.Close()
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Error reading message: %v\n", err)
			continue
		}

		log.Printf("Consumed message: %s", string(msg.Value))

		if err := handler(msg.Key, msg.Value); err != nil {
			log.Printf("Handler error: %v\n", err)
		}
	}
}
