package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis"
)

var ctx = context.Background()
var rdb *redis.Client

func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password
		DB:       0,  // use default DB
	})
}

type User struct {
	ID       string
	Usernama string
	Email    string
}

func createSession(user User, sessionToken string) error {
	key := fmt.Sprintf("session:%s", sessionToken)
	err := rdb.HMSet(key, map[string]interface{}{
		"id":       user.ID,
		"username": user.Usernama,
		"email":    user.Email,
	}).Err()
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func getSession(sessionToken string) (User, error) {
	key := fmt.Sprintf("session:%s", sessionToken)
	val, err := rdb.HGetAll(key).Result()
	if err != nil {
		return User{}, err
	}

	return User{
		ID:       val["id"],
		Usernama: val["username"],
		Email:    val["email"],
	}, nil
}

func deleteSession(sessionToken string) error {
	key := fmt.Sprintf("session:%s", sessionToken)
	err := rdb.Del(key).Err()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	initRedis()
	sessionToken := "1234567890"

	// Endpoint login
	user := User{
		ID:       "1",
		Usernama: "johndoe",
		Email:    "asdasd@mail.com",
	}
	createSession(user, sessionToken)
	// Endpoint login

	// Endpoint dashboard
	userSession, err := getSession(sessionToken)
	if err != nil {
		fmt.Println(err)
	}
	if userSession.ID == "" {
		fmt.Println("Session not found. Please login first")
	} else {
		fmt.Println(userSession)
	}
	// Endpoint dashboard

	// Endpoint logout
	deleteSession(sessionToken)
	// Endpoint logout
}
