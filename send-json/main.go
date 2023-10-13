package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type Message struct {
	Text  string `json:"text"`
	Count int    `json:"count"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://juscilan:pjuscilan@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // Queue name
		false,   // Durable
		false,   // Delete when unused
		false,   // Exclusive
		false,   // No-wait
		nil,     // Arguments
	)
	failOnError(err, "Failed to declare a queue")

	message := Message{
		Text:  "Hello, RabbitMQ!, This message has been produced by juscilan.com",
		Count: 42,
	}

	jsonMessage, err := json.Marshal(message)
	failOnError(err, "Failed to marshal JSON message")

	err = ch.Publish(
		"",     // Exchange
		q.Name, // Routing key
		false,  // Mandatory
		false,  // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(jsonMessage),
		},
	)
	failOnError(err, "Failed to publish a JSON message")

	fmt.Printf("Sent JSON message: %s\n", jsonMessage)
}
