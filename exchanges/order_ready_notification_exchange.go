package exchanges

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func DeclareOrderReadyNotificationExchange(ch *amqp.Channel) {
	err := ch.ExchangeDeclare(
		"order_ready_exchange",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare order_ready_exchange: %v", err)
	}
}
