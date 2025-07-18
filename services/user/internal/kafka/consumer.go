package kafka

import (
	"context"
	"encoding/json"
	"github.com/mdqni/Attendly/services/user/internal/service"
	"github.com/mdqni/Attendly/shared/domain"
	"github.com/mdqni/Attendly/shared/kafka"
	"log"
)

type EventConsumer struct {
	svc   service.UserService
	kafka *kafka.Consumer
}

func NewEventConsumer(brokers string, svc service.UserService) (*EventConsumer, error) {
	k, err := kafka.NewConsumer(brokers, "user-consumer-group", []string{"auth.user_registered"})
	if err != nil {
		return nil, err
	}
	return &EventConsumer{
		svc:   svc,
		kafka: k,
	}, nil
}

func (e *EventConsumer) Start(ctx context.Context) error {
	return e.kafka.Start(ctx, func(key, value []byte) error {
		var event domain.UserRegisteredEvent
		if err := json.Unmarshal(value, &event); err != nil {
			log.Printf("Failed to unmarshal event: %v", err)
			return err
		}
		log.Printf("Received user_registered: %+v", event)

		user := domain.User{
			ID:    event.UserID,
			Email: event.Email,
			Role:  event.Role,
		}

		result, err := e.svc.CreateUser(ctx, &user)
		if err != nil {
			log.Printf("Failed to save user: %v", err)
		}
		log.Printf("Saved user: %v", result)
		return nil
	})
}
