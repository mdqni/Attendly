package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:  brokers,
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		}),
	}
}

func (p *Producer) Produce(key, value []byte) error {
	msg := kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}

	log.Printf("Producing message: %s", value)
	return p.writer.WriteMessages(context.Background(), msg)
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
