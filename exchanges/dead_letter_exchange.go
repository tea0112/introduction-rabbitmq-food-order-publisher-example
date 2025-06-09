package exchanges

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func DeclareDeadLetterExchange(ch *amqp.Channel) {
	err := ch.ExchangeDeclare(
		"dlx_exchange",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare dlx_exchange: %v", err)
	}
}
