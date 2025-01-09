package main

import (
	"fmt"
	"log"
	"time"

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
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// exchangeDeclare
	err = ch.ExchangeDeclare(
		"transactionUser", // name
		"fanout",          // type
		true,              // durable
		false,             // auto-deleted
		false,             // internal
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare an exchange: %v", err)
	}

	for i := 0; i < 10; i++ {
		datenow := time.Now()
		body := "Hello World! " + datenow.Format("2006-01-02 15:04:05")
		err = ch.Publish(
			"transactionUser", // exchange
			"",                // routing key
			false,             // mandatory
			false,             // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		if err != nil {
			log.Fatalf("Failed to publish a message: %v", err)
		}
		time.Sleep(1 * time.Second)
	}

	fmt.Println("Message Published")
}
