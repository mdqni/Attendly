package kafka

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log"
)

type MessageHandler func(key, value []byte) error

type Consumer struct {
	c *kafka.Consumer
}

func NewConsumer(brokers, groupID string, topics []string) (*Consumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	err = c.SubscribeTopics(topics, nil)
	if err != nil {
		return nil, err
	}

	return &Consumer{c: c}, nil
}

func (c *Consumer) Start(ctx context.Context, handler MessageHandler) error {
	for {
		select {
		case <-ctx.Done():
			log.Println("Consumer shutting down...")
			return c.c.Close()
		default:
			msg, err := c.c.ReadMessage(-1)
			if err != nil {
				log.Printf("Error reading message: %v\n", err)
				continue
			}

			if err := handler(msg.Key, msg.Value); err != nil {
				log.Printf("Handler error: %v\n", err)
			}
		}
	}
}
