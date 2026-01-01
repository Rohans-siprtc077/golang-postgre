package config

import (
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var RabbitConn *amqp.Connection
var RabbitChannel *amqp.Channel

func mustEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Missing required environment variable: %s", key)
	}
	return val
}

func ConnectRabbitMQ() {

	url := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		mustEnv("RABBITMQ_USER"),
		mustEnv("RABBITMQ_PASSWORD"),
		mustEnv("RABBITMQ_HOST"),
		mustEnv("RABBITMQ_PORT"),
	)

	var err error

	for i := 1; i <= 10; i++ {
		RabbitConn, err = amqp.Dial(url)
		if err == nil {
			log.Println(" RabbitMQ connected")
			break
		}

		log.Printf("RabbitMQ not ready (attempt %d/10). Retrying...", i)
		time.Sleep(3 * time.Second)
	}

	if RabbitConn == nil {
		log.Fatal("RabbitMQ connection failed after retries")
	}

	// Create channel
	RabbitChannel, err = RabbitConn.Channel()
	if err != nil {
		log.Fatal("RabbitMQ channel creation failed:", err)
	}

	// Read configuration
	exchange := mustEnv("RABBITMQ_EXCHANGE")
	dlx := mustEnv("RABBITMQ_DLX")
	createdQueue := mustEnv("RABBITMQ_CREATED_QUEUE")
	updatedQueue := mustEnv("RABBITMQ_UPDATED_QUEUE")

	// Declare main exchange
	if err := RabbitChannel.ExchangeDeclare(
		exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		log.Fatal("Failed to declare exchange:", err)
	}

	// Declare dead-letter exchange
	if err := RabbitChannel.ExchangeDeclare(
		dlx,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		log.Fatal("Failed to declare DLX:", err)
	}

	// Queue arguments (DLX)
	args := amqp.Table{
		"x-dead-letter-exchange": dlx,
	}

	// Declare queues
	if _, err := RabbitChannel.QueueDeclare(
		createdQueue,
		true,
		false,
		false,
		false,
		args,
	); err != nil {
		log.Fatal("Failed to declare created queue:", err)
	}

	if _, err := RabbitChannel.QueueDeclare(
		updatedQueue,
		true,
		false,
		false,
		false,
		args,
	); err != nil {
		log.Fatal("Failed to declare updated queue:", err)
	}

	// Bind queues
	if err := RabbitChannel.QueueBind(
		createdQueue,
		"user.created",
		exchange,
		false,
		nil,
	); err != nil {
		log.Fatal("Failed to bind created queue:", err)
	}

	if err := RabbitChannel.QueueBind(
		updatedQueue,
		"user.updated",
		exchange,
		false,
		nil,
	); err != nil {
		log.Fatal("Failed to bind updated queue:", err)
	}

	log.Println("RabbitMQ setup completed successfully")
}
