package kafka

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/mdqni/Attendly/services/user/internal/service"
	"github.com/mdqni/Attendly/shared/domain"
	sharedkafka "github.com/mdqni/Attendly/shared/kafka"
)

type EventConsumer struct {
	svc      service.UserService
	consumer *sharedkafka.Consumer
}

func NewEventConsumer(brokersCSV, topic, groupID string, svc service.UserService) (*EventConsumer, error) {
	brokers := strings.Split(brokersCSV, ",")
	c := sharedkafka.NewConsumer(brokers, topic, groupID)

	return &EventConsumer{
		svc:      svc,
		consumer: c,
	}, nil
}

func (e *EventConsumer) Start(ctx context.Context) error {
	return e.consumer.Start(ctx, func(key, value []byte) error {
		var event domain.UserRegisteredEvent
		if err := json.Unmarshal(value, &event); err != nil {
			log.Printf("Failed to unmarshal event: %v", err)
			return err
		}

		user := domain.User{
			ID:    event.UserID,
			Email: event.Email,
			Role:  event.Role,
			Name:  event.Name,
		}

		if _, err := e.svc.CreateUser(ctx, &user); err != nil {
			log.Printf("Failed to save user: %v", err)
		}
		return nil
	})
}
