package publishes

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

func PublishOrderReadyNotification(ch *amqp.Channel, batchNum int) {
	// 2. Fanout Exchange: Order Ready Notification
	ready := map[string]string{
		"order_id": fmt.Sprintf("order_%d", batchNum),
	}
	readyJSON, _ := json.Marshal(ready)
	err := ch.PublishWithContext(
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
}
