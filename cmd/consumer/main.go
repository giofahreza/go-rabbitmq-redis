package main

import (
	"encoding/json"
	"fmt"
	"log"

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

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func createSession(user User, sessionID string) error {
	// Create a new session
	key := fmt.Sprintf("session:%s", sessionID)
	err := rdb.HMSet(key, map[string]interface{}{
		"name":     user.Name,
		"email":    user.Email,
		"username": user.Username,
	}).Err()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

func getSession(sessionID string) (User, error) {
	// Get the session
	key := fmt.Sprintf("session:%s", sessionID)
	val, err := rdb.HGetAll(key).Result()
	if err != nil {
		return User{}, fmt.Errorf("failed to get session: %w", err)
	}

	// Unmarshal the session data
	var user User
	userData, err := json.Marshal(val)
	if err != nil {
		return User{}, fmt.Errorf("failed to marshal session data: %w", err)
	}
	err = json.Unmarshal(userData, &user)
	if err != nil {
		return User{}, fmt.Errorf("failed to unmarshal session data: %w", err)
	}
	return user, nil
}

func deleteSession(sessionID string) error {
	// Delete the session
	key := fmt.Sprintf("session:%s", sessionID)
	err := rdb.Del(key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}

func main() {
	initRedis()

	val, err := rdb.HGet("message", "message").Result()
	if err != nil {
		log.Fatalf("Failed to get message from Redis: %s", err)
	}
	log.Printf("Message from Redis: %s", val)

	token := "123qweasdzxc12sda8yeh9ui3"
	user := User{
		Name:     "John Doe",
		Email:    "john@mail.com",
		Username: "johndoe",
	}
	err = createSession(user, token)
	if err != nil {
		log.Fatalf("Failed to create session: %s", err)
	}
	log.Printf("Session created for user: %s", user.Username)

	// Get the session
	retrievedUser, err := getSession(token)
	if err != nil {
		log.Fatalf("Failed to get session: %s", err)
	}
	log.Printf("Retrieved session for user: %s", retrievedUser.Username)
	// Print the session data
	log.Printf("Session data: %+v", retrievedUser)

	// Delete the session
	err = deleteSession(token)
	if err != nil {
		log.Fatalf("Failed to delete session: %s", err)
	}
	log.Printf("Session deleted for user: %s", user.Username)

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
		log.Fatalf("Failed to register a consumer: %s", err)
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf(" [x] Received %s", d.Body)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
