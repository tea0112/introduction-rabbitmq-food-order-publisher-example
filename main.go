package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
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

	// Enable publisher confirms
	err = ch.Confirm(false)
	if err != nil {
		log.Fatalf("Failed to enable publisher confirms: %v", err)
	}

	// Set up confirmation channel
	confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 10))

	// Start a goroutine to handle confirmations
	go func() {
		for confirm := range confirms {
			if confirm.Ack {
				log.Printf("Message %d confirmed", confirm.DeliveryTag)
			} else {
				log.Printf("Message %d not confirmed", confirm.DeliveryTag)
			}
		}
	}()

	// Declare Direct Exchange for order placement
	err = ch.ExchangeDeclare(
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

	// Declare Fanout Exchange for order ready notifications
	err = ch.ExchangeDeclare(
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

	// Declare Topic Exchange for delivery assignments
	err = ch.ExchangeDeclare(
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

	// Declare Headers Exchange for status updates
	err = ch.ExchangeDeclare(
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

	// Declare Dead Letter Exchange
	err = ch.ExchangeDeclare(
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

	// Simulate publishing messages
	for i := 0; i < 5; i++ {
		// 1. Direct Exchange: Order Placement (including a poison message)
		order := map[string]string{
			"order_id":      fmt.Sprintf("order_%d", i),
			"restaurant_id": "restaurant_abc",
		}
		if i == 2 { // Inject a poison message
			order["order_id"] = "order_id_POISON"
		}
		orderJSON, _ := json.Marshal(order)
		err = ch.PublishWithContext(
			context.Background(), // context
			"orders_exchange",    // exchange
			"restaurant_abc",     // routing key
			false,                // mandatory
			false,                // immediate
			amqp.Publishing{
				ContentType: "application/json", // content type
				Body:        orderJSON,          // message body
			},
		)
		if err != nil {
			log.Printf("Failed to publish order: %v", err)
		} else {
			log.Printf("Published order %s to restaurant_abc", order["order_id"])
		}
		time.Sleep(time.Millisecond * 300) // Simulate delay between publishes

		// 2. Fanout Exchange: Order Ready Notification
		ready := map[string]string{
			"order_id": fmt.Sprintf("order_%d", i),
		}
		readyJSON, _ := json.Marshal(ready)
		err = ch.PublishWithContext(
			context.Background(),
			"order_ready_exchange",
			"", // no routing key for fanout
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        readyJSON,
			},
		)
		if err != nil {
			log.Printf("Failed to publish order ready: %v", err)
		} else {
			log.Printf("Published order %s ready notification", ready["order_id"])
		}
		time.Sleep(time.Millisecond * 300) // Simulate delay between publishes

		// 3. Topic Exchange: Delivery Assignment
		assignment := map[string]string{
			"order_id":  fmt.Sprintf("order_%d", i),
			"driver_id": "driver_xyz",
			"region":    "north",
		}
		assignmentJSON, _ := json.Marshal(assignment)
		err = ch.PublishWithContext(
			context.Background(),
			"delivery_exchange",
			"delivery.assign.north", // routing key
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        assignmentJSON,
			},
		)
		if err != nil {
			log.Printf("Failed to publish assignment: %v", err)
		} else {
			log.Printf("Published delivery assignment for order %s", assignment["order_id"])
		}
		time.Sleep(time.Millisecond * 300) // Simulate delay between publishes

		// 4. Headers Exchange: Status Update with TTL
		status := map[string]string{
			"order_id": fmt.Sprintf("order_%d", i),
			"status":   "in_transit",
		}
		statusJSON, _ := json.Marshal(status)
		err = ch.PublishWithContext(
			context.Background(),
			"status_exchange",
			"", // no routing key for headers
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        statusJSON,
				Headers:     amqp.Table{"status": "in_transit"}, // headers for routing
				Expiration:  "1",                                // TTL of 5 seconds
			},
		)
		if err != nil {
			log.Printf("Failed to publish status: %v", err)
		} else {
			log.Printf("Published status update for order %s", status["order_id"])
		}

		time.Sleep(time.Second) // Delay between iterations
	}
}
