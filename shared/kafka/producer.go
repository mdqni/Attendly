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
	deliveryChan := make(chan kafka.Event)

	err := p.p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            key,
		Value:          value,
	}, deliveryChan)

	if err != nil {
		return err
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		log.Printf("Delivery failed: %v", m.TopicPartition.Error)
		return m.TopicPartition.Error
	}

	log.Printf("Delivered message to %v", m.TopicPartition)
	close(deliveryChan)

	return nil
}

func (p *Producer) Close() {
	p.p.Flush(5000)
	p.p.Close()
}
