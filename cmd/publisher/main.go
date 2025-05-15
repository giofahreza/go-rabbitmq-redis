package main

import (
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/streadway/amqp"
)

var (
	rdb *redis.Client
)

func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := rdb.Ping().Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %s", err)
	}
	log.Println("Connected to Redis")
}

func main() {
	initRedis()

	err := rdb.HMSet("message", map[string]interface{}{"message": "Something"}).Err()
	if err != nil {
		log.Fatalf("Failed to set message in Redis: %s", err)
	}

	conn, err := amqp.Dial("amqps://otejfslf:vm7ZCp3RBW10r5Q8drv088o4JNjDIB09@fuji.lmq.cloudamqp.com/otejfslf")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}

	datetimenow := time.Now()
	body := "Hello World! " + datetimenow.Format("2006-01-02 15:04:05")
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	if err != nil {
		log.Fatalf("Failed to publish a message: %s", err)
	}
	log.Printf(" [x] Sent %s", body)
}
