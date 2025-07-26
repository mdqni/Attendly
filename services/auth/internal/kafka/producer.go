package kafka

import (
	"context"
	"encoding/json"
	"github.com/mdqni/Attendly/shared/domain"
	"github.com/mdqni/Attendly/shared/kafka"
	"log"
)

type EventProducer struct {
	producer *kafka.Producer
}

func NewEventProducer(brokers []string, topic string) *EventProducer {
	return &EventProducer{
		producer: kafka.NewProducer(brokers, topic),
	}
}

func (e *EventProducer) SendUserRegisteredEvent(ctx context.Context, userID, email, role, name string) error {
	event := domain.UserRegisteredEvent{
		UserID: userID,
		Name:   name,
		Email:  email,
		Role:   role,
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	log.Printf("Sending Kafka Event: %s", data)
	return e.producer.Produce([]byte(userID), data)
}

func (e *EventProducer) Close() {
	_ = e.producer.Close()
}
