package exchanges

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func DeclareDeliveryAssignmentExchange(ch *amqp.Channel) {
	err := ch.ExchangeDeclare(
		"delivery_exchange",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare delivery_exchange: %v", err)
	}
}
