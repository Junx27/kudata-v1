package user

import (
	"fmt"
	"log"
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

func (ue *UserEvent) SubscribeUser() {
	q, err := ue.Channel.QueueDeclare(
		"user_created_info", // random queue name
		true,                // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		log.Fatalf("%s", err)
	}

	err = ue.Channel.QueueBind(
		q.Name,              // queue name
		"user.created.info", // routing key
		event.ExchangeName,  // exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("%s", err)
	}

	msgs, err := ue.Channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("%s", err)
	}

	for msg := range msgs {
		ue.handleConsumeSomething(msg)
	}
}

func (ue *UserEvent) handleConsumeSomething(msg amqp.Delivery) {
	message := string(msg.Body)
	fmt.Printf("Received message from survey service: %s\n", message)
}
