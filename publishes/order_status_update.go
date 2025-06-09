package publishes

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

func PublishOrderUpdateStatus(ch *amqp.Channel, batchNum int) {
	// 4. Headers Exchange: Status Update with TTL
	status := map[string]string{
		"order_id": fmt.Sprintf("order_%d", batchNum),
		"status":   "in_transit",
	}
	statusJSON, _ := json.Marshal(status)
	err := ch.PublishWithContext(
		context.Background(),
		"status_exchange",
		"", // no routing key for headers
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        statusJSON,
			Headers:     amqp.Table{"status": "in_transit"}, // headers for routing
			Expiration:  "500",                              // TTL of 5 seconds
		},
	)
	if err != nil {
		log.Printf("Failed to publish status: %v", err)
	} else {
		log.Printf("Published status update for order %s", status["order_id"])
	}

	time.Sleep(time.Second) // Delay between iterations
}
