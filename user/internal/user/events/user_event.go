package events

import (
	"context"
	"encoding/json"
	"log"
	"user/internal/user/model"
	"user/internal/user/repository"
	"user/pkg/event"

	"github.com/streadway/amqp"
)

type UserEvent struct {
	Channel *amqp.Channel
}

func NewUserEvent(ch *amqp.Channel) UserEvent {
	return UserEvent{
		Channel: ch,
	}
}

func (ue *UserEvent) SubscribeCreateUser() {
	q, err := ue.Channel.QueueDeclare(
		"create_user", // queue name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		log.Fatalf("QueueDeclare error: %v", err)
	}

	err = ue.Channel.QueueBind(
		q.Name,             // queue name
		"create.user",      // routing key
		event.ExchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("QueueBind error: %v", err)
	}

	msgs, err := ue.Channel.Consume(
		q.Name, // queue
		"",     // consumer tag
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Consume error: %v", err)
	}

	for msg := range msgs {
		go ue.handleConsumeCreateUser(msg)
	}
}

func (ue *UserEvent) handleConsumeCreateUser(msg amqp.Delivery) {

	var payload model.UserEvent
	if err := json.Unmarshal(msg.Body, &payload); err != nil {
		log.Printf("Error parsing message: %v", err)
		return
	}

	user := model.UserInput{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: payload.Password,
	}
	if err := repository.StoreUser(context.Background(), user); err != nil {
		log.Printf("Error storing survey: %v", err)
		return
	}

	log.Printf("✅ Create user successfully")
}
