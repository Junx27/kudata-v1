package events

import (
	"context"
	"encoding/json"
	"log"
	"path/filepath"
	"survey/internal/survey/model"
	"survey/internal/survey/repository"
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

func (se *SurveyEvent) SubscribeCreateSurvey() {
	q, err := se.Channel.QueueDeclare(
		"create_survey", // queue name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		log.Fatalf("QueueDeclare error: %v", err)
	}

	err = se.Channel.QueueBind(
		q.Name,             // queue name
		"create.survey",    // routing key
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
		go se.handleConsumeCreateSurvey(msg)
	}
}

func (se *SurveyEvent) handleConsumeCreateSurvey(msg amqp.Delivery) {

	var payload model.SurveyEvent
	if err := json.Unmarshal(msg.Body, &payload); err != nil {
		log.Printf("Error parsing message: %v", err)
		return
	}

	fileName := filepath.Base(payload.Image)
	survey := model.SurveyInput{
		Name:        payload.Name,
		Price:       payload.Price,
		Description: payload.Description,
		CategoryID:  payload.CategoryID,
	}
	if err := repository.StoreSurvey(context.Background(), survey, fileName); err != nil {
		log.Printf("Error storing survey: %v", err)
		return
	}

	log.Printf("✅ Create successfully")
}
func (se *SurveyEvent) SubscribeUpdateSurvey() {
	q, err := se.Channel.QueueDeclare(
		"update_survey", // queue name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		log.Fatalf("QueueDeclare error: %v", err)
	}

	err = se.Channel.QueueBind(
		q.Name,             // queue name
		"update.survey",    // routing key
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
		go se.handleConsumeUpdateSurvey(msg)
	}
}

func (se *SurveyEvent) handleConsumeUpdateSurvey(msg amqp.Delivery) {

	var payload model.SurveyEvent
	if err := json.Unmarshal(msg.Body, &payload); err != nil {
		log.Printf("Error parsing message: %v", err)
		return
	}

	fileName := filepath.Base(payload.Image)
	id := payload.ID
	survey := model.SurveyInput{
		Name:        payload.Name,
		Price:       payload.Price,
		Description: payload.Description,
		CategoryID:  payload.CategoryID,
	}
	if err := repository.UpdateSurvey(context.Background(), id, survey, fileName); err != nil {
		log.Printf("Error storing survey: %v", err)
		return
	}

	log.Printf("✅ Update successfully")
}
func (se *SurveyEvent) SubscribeDeleteSurvey() {
	q, err := se.Channel.QueueDeclare(
		"delete_survey", // queue name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		log.Fatalf("QueueDeclare error: %v", err)
	}

	err = se.Channel.QueueBind(
		q.Name,             // queue name
		"delete.survey",    // routing key
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
		go se.handleConsumeDeleteSurvey(msg)
	}
}

func (se *SurveyEvent) handleConsumeDeleteSurvey(msg amqp.Delivery) {

	var payload model.SurveyEvent
	if err := json.Unmarshal(msg.Body, &payload); err != nil {
		log.Printf("Error parsing message: %v", err)
		return
	}
	id := payload.ID
	if err := repository.DeleteSurvey(context.Background(), id); err != nil {
		log.Printf("Error storing survey: %v", err)
		return
	}

	log.Printf("✅ Delete successfully")
}
