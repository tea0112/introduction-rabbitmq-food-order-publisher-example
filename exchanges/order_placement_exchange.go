package exchanges

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func DeclareOrderPlacementExchange(ch *amqp.Channel) {
	// Declare Direct Exchange for order placement
	err := ch.ExchangeDeclare(
		"orders_exchange", // name
		"direct",          // type
		true,              // durable
		false,             // auto-deleted
		false,             // internal
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare orders_exchange: %v", err)
	}
}
