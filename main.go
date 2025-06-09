package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"introduction-rabbitmq-food-order-publisher-example/exchanges"
	"introduction-rabbitmq-food-order-publisher-example/publishes"
	"log"
)

func main() {
	// Connect to RabbitMQ server
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Open a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	log.Println("Starting to publish messages 🚀")

	// Enable publisher confirms
	err = ch.Confirm(false)
	if err != nil {
		log.Fatalf("Failed to enable publisher confirms: %v", err)
	}

	// Set up confirmation channel
	confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 10))

	// Start a goroutine to handle confirmations
	go handlePublishConfirmations(confirms)

	// Declare Direct Exchange for order placement
	exchanges.DeclareOrderPlacementExchange(ch)

	// Declare Fanout Exchange for order ready notifications
	exchanges.DeclareOrderReadyNotificationExchange(ch)

	// Declare Topic Exchange for delivery assignments
	exchanges.DeclareDeliveryAssignmentExchange(ch)

	// Declare Headers Exchange for status updates
	exchanges.DeclareOrderStatusUpdateExchange(ch)

	// Declare Dead Letter Exchange
	exchanges.DeclareDeadLetterExchange(ch)

	// Simulate publishing messages
	for batchNum := 0; batchNum < 5; batchNum++ {
		// 1. Direct Exchange: Order Placement (including a poison message)
		publishes.PublishOrderPlacement(ch, batchNum)

		// 2. Fanout Exchange: Order Ready Notification
		publishes.PublishOrderReadyNotification(ch, batchNum)

		// 3. Topic Exchange: Delivery Assignment
		publishes.PublishDeliveryAssignment(ch, batchNum)

		// 4. Headers Exchange: Status Update with TTL
		publishes.PublishOrderUpdateStatus(ch, batchNum)
	}
}

func handlePublishConfirmations(confirms chan amqp.Confirmation) {
	for confirm := range confirms {
		if confirm.Ack {
			log.Printf("Message %d confirmed", confirm.DeliveryTag)
		} else {
			log.Printf("Message %d not confirmed", confirm.DeliveryTag)
		}
	}
}
