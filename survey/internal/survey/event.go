package survey

import (
	"fmt"
	"log"
	"survey/pkg/event"

	"github.com/streadway/amqp"
)

type SurveyEvent struct {
	Channel *amqp.Channel
}

func NewSurveyEvent(ch *amqp.Channel) SurveyEvent {
	return SurveyEvent{
		Channel: ch,
	}
}

func (se *SurveyEvent) SubscribeSurvey() {
	q, err := se.Channel.QueueDeclare(
		"user_created_success", // random queue name
		true,                   // durable
		false,                  // delete when unused
		false,                  // exclusive
		false,                  // no-wait
		nil,                    // arguments
	)
	if err != nil {
		log.Fatalf("%s", err)
	}

	err = se.Channel.QueueBind(
		q.Name,                // queue name
		"create.user.success", // routing key
		event.ExchangeName,    // exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("%s", err)
	}

	msgs, err := se.Channel.Consume(
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
		se.handleConsumeSomething(msg)
	}
}

func (se *SurveyEvent) handleConsumeSomething(msg amqp.Delivery) {
	message := string(msg.Body)
	fmt.Printf("Received message from user service: %s\n", message)
}
