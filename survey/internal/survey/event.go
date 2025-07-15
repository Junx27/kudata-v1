package survey

import (
	"fmt"
	"log"
	"survey/pkg/event"
	"sync"

	"github.com/streadway/amqp"
)

type SurveyEvent struct {
	Channel *amqp.Channel
}

var (
	LatestMessage string
	MessageMu     sync.RWMutex
)

func NewSurveyEvent(ch *amqp.Channel) SurveyEvent {
	return SurveyEvent{
		Channel: ch,
	}
}

// Fungsi ini sekarang menerima parameter: broadcast ke WebSocket
func (se *SurveyEvent) SubscribeSurvey(broadcastFunc func(string)) {
	q, err := se.Channel.QueueDeclare(
		"user_created_success",
		true, false, false, false, nil,
	)
	if err != nil {
		log.Fatalf("%s", err)
	}

	err = se.Channel.QueueBind(
		q.Name,
		"create.user.success",
		event.ExchangeName,
		false, nil,
	)
	if err != nil {
		log.Fatalf("%s", err)
	}

	msgs, err := se.Channel.Consume(
		q.Name, "", true, false, false, false, nil,
	)
	if err != nil {
		log.Fatalf("%s", err)
	}

	for msg := range msgs {
		se.handleConsumeSomething(msg, broadcastFunc)
	}
}

func (se *SurveyEvent) handleConsumeSomething(msg amqp.Delivery, broadcastFunc func(string)) {
	message := string(msg.Body)

	MessageMu.Lock()
	LatestMessage = message
	MessageMu.Unlock()

	fmt.Printf("Received message from user service: %s\n", message)

	// Kirim ke WebSocket
	broadcastFunc(message)
}
