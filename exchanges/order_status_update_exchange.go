package exchanges

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func DeclareOrderStatusUpdateExchange(ch *amqp.Channel) {
	err := ch.ExchangeDeclare(
		"status_exchange",
		"headers",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare status_exchange: %v", err)
	}
}
