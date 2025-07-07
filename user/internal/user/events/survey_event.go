package events

import (
	"encoding/json"
	"fmt"
	"log"
	"user/pkg/event"

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
		"res_all_survey", // queue name
		true,             // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		log.Fatalf("QueueDeclare error: %v", err)
	}

	err = se.Channel.QueueBind(
		q.Name,             // queue name
		"res.all.survey",   // routing key
		event.ExchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("QueueBind error: %v", err)
	}

	msgs, err := se.Channel.Consume(
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
		go se.handleConsumeResponse(msg)
	}
}

func (se *SurveyEvent) handleConsumeResponse(msg amqp.Delivery) {
	fmt.Println("========== Received Response from Survey Service ==========")

	var surveys []map[string]interface{}
	if err := json.Unmarshal(msg.Body, &surveys); err != nil {
		log.Printf("Failed to parse survey response as JSON: %v", err)
		fmt.Printf("Raw message: %s\n", string(msg.Body))
		return
	}

	for i, s := range surveys {
		fmt.Printf("Survey #%d:\n", i+1)
		for k, v := range s {
			fmt.Printf("  %s: %v\n", k, v)
		}
	}
	fmt.Println("==========================================================")
}
