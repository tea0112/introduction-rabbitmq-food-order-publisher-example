package publishes

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

func PublishOrderPlacement(ch *amqp.Channel, batchNum int) {
	// 1. Direct Exchange: Order Placement (including a poison message)
	orderPlacementBody := map[string]string{
		"order_id":      fmt.Sprintf("order_%d", batchNum),
		"restaurant_id": "restaurant_abc",
	}
	if batchNum == 2 { // Inject a poison message
		orderPlacementBody["order_id"] = "order_id_POISON"
	}
	orderJSON, _ := json.Marshal(orderPlacementBody)
	err := ch.PublishWithContext(
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
		log.Printf("Published order %s to restaurant_abc", orderPlacementBody["order_id"])
	}

	time.Sleep(time.Millisecond * 300) // Simulate delay between publishes
}
