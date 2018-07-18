package main

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

const QueueName = "orm:todos"

func Connect2AMQPAndSetupQueue(amqpUri string, todoChanged chan Todo) {
	connection, err := amqp.Dial(amqpUri)
	if err != nil {
		log.Fatalf("Failed to connect to %s: %s", amqpUri, err)
	}
	defer connection.Close()
	log.Printf("Connected to AMQP: %s", amqpUri)

	channel, err := connection.Channel()
	if err != nil {
		log.Fatalf("Failed to acquire channel: %s", err)
	}
	defer channel.Close()

	_, err = channel.QueueDeclare(QueueName, false, false, false, false, nil, )
	if err != nil {
		log.Fatalf("Failed to declare queue '%s': %s", QueueName, err)
	}
	log.Printf("Using queue: %s", QueueName)

	for {
		todo := <-todoChanged
		todoBytes, _ := json.Marshal(todo)

		channel.Publish("", QueueName, false, false,
			amqp.Publishing{ContentType: "text/json", Body: []byte(todoBytes)},
		)
	}

}
