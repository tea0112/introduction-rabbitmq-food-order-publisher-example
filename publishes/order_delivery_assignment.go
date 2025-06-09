package publishes

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

func PublishDeliveryAssignment(ch *amqp.Channel, batchNum int) {
	// 3. Topic Exchange: Delivery Assignment
	assignment := map[string]string{
		"order_id":  fmt.Sprintf("order_%d", batchNum),
		"driver_id": "driver_xyz",
		"region":    "north",
	}
	assignmentJSON, _ := json.Marshal(assignment)
	err := ch.PublishWithContext(
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
}
