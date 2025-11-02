package consumers

import (
	"log"
	"order-service/internal/rabbitmq"

	"github.com/streadway/amqp"
)

func StartOrderLoggerConsumer(rabbitmqURL string) {
	log.Println("Starting RabbitMQ consumer for logging...")

	eventHandler := func(d amqp.Delivery) {
		log.Printf("Successfully processed event from queue: %s", string(d.Body))
	}

	if err := rabbitmq.SetupConsumer(rabbitmqURL, "amq.topic", "order.created", "order_queue", eventHandler); err != nil {
		log.Fatalf("Failed to start RabbitMQ consumer for logging: %v", err)
	}
}

func StartProductEventListener(rabbitmqURL string) {
	log.Println("Starting RabbitMQ consumer for product events...")

	productEventHandler := func(d amqp.Delivery) {
		log.Printf("Received a product event: %s", string(d.Body))
	}

	if err := rabbitmq.SetupConsumer(rabbitmqURL, "amq.topic", "product.created", "product_queue", productEventHandler); err != nil {
		log.Fatalf("Failed to start product event consumer: %v", err)
	}
}
