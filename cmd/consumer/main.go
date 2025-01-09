package main

import (
	"log"

	"github.com/streadway/amqp"
)

func main() {
	conn, err := amqp.Dial("amqps://vdfkrprg:rlM-PCmX1b1ulLdlQOaueBZeNxYSzokd@fuji.lmq.cloudamqp.com/vdfkrprg")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open channel: %v", err)
	}
	defer ch.Close()

	// Declare the same exchange as publisher
	exchangeName := "transactionUser"

	// Declare a queue with a random name
	q, err := ch.QueueDeclare(
		"",    // name (empty = let RabbitMQ generate a name)
		false, // durable
		true,  // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	err = ch.QueueBind(
		q.Name,       // queue name
		"",           // routing key
		exchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	log.Printf("Consumer started. Queue: %s\n", q.Name)

	forever := make(chan bool)
	go func() {
		for msg := range msgs {
			log.Printf("Received: %s\n", msg.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
