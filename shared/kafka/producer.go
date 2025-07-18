package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log"
)

type Producer struct {
	p *kafka.Producer
}

func NewProducer(brokers string) (*Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
	})
	if err != nil {
		return nil, err
	}
	return &Producer{p: p}, nil
}

func (p *Producer) Produce(topic string, key, value []byte) error {
	err := p.p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            key,
		Value:          value,
	}, nil)
	if err != nil {
		return err
	}
	go func() {
		for e := range p.p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("Delivery failed: %v\n", ev.TopicPartition.Error)
				} else {
					log.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()
	return nil
}

func (p *Producer) Close() {
	p.p.Close()
}
