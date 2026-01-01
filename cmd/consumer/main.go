package main

import (
	"encoding/json"
	"log"

	"golang-postgre/config"
	"golang-postgre/events"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	config.ConnectRabbitMQ()

	createdQueue := "user.created.queue"
	updatedQueue := "user.updated.queue"

	createdMsgs, err := config.RabbitChannel.Consume(
		createdQueue,
		"created-consumer",
		false, // manual ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("❌ Failed to consume created queue:", err)
	}

	updatedMsgs, err := config.RabbitChannel.Consume(
		updatedQueue,
		"updated-consumer",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("❌ Failed to consume updated queue:", err)
	}

	log.Println("📥 Consumer started. Waiting for events...")

	forever := make(chan bool)

	go func() {
		for msg := range createdMsgs {
			var event events.UserEvent

			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Println("❌ Invalid message:", err)
				msg.Nack(false, false)
				continue
			}

			log.Printf("[USER_CREATED] Welcome email sent to %s\n", event.Data.Email)
			msg.Ack(false)
		}
	}()

	go func() {
		for msg := range updatedMsgs {
			var event events.UserEvent

			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Println("❌ Invalid message:", err)
				msg.Nack(false, false)
				continue
			}

			log.Printf("[USER_UPDATED] User %d profile updated\n", event.Data.UserID)
			msg.Ack(false)
		}
	}()

	<-forever
}
