package kafka

import (
	"context"
	"encoding/json"
	"github.com/mdqni/Attendly/shared/domain"
	"github.com/mdqni/Attendly/shared/kafka"
	"log"
)

type EventProducer struct {
	kafkaProducer *kafka.Producer
	topic         string
}

func NewEventProducer(kafkaBrokers string) (*EventProducer, error) {
	p, err := kafka.NewProducer(kafkaBrokers)
	if err != nil {
		return nil, err
	}
	return &EventProducer{
		kafkaProducer: p,
		topic:         "auth.user_registered",
	}, nil
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

	log.Printf("Sending kafka to Kafka: %s", data)
	return e.kafkaProducer.Produce(e.topic, []byte(userID), data)
}

func (e *EventProducer) Close() {
	e.kafkaProducer.Close()
}
