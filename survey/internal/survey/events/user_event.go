package events

import (
	"log"
	"survey/pkg/event"

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
		"req_all_survey", // queue name
		true,             // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	err = ue.Channel.QueueBind(
		q.Name,             // queue name
		"req.all.survey",   // routing key
		event.ExchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to bind queue: %v", err)
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
		log.Fatalf("Failed to register consumer: %v", err)
	}

	for msg := range msgs {
		go ue.handleConsumeSomething(msg)
	}
}

func (ue *UserEvent) handleConsumeSomething(msg amqp.Delivery) {
	log.Printf("Received message from user service: %s", string(msg.Body))
}
