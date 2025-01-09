package main

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis"
)

var ctx = context.Background()

type User struct {
	ID   string
	Name string
}

func createClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password
		DB:       0,  // use default DB
	})
	return client
}

func setKey(client *redis.Client, key, value string) error {
	err := client.Set(key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func getKey(client *redis.Client, key string) (string, error) {
	val, err := client.Get(key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func createUser(client *redis.Client, user User) error {
	err := client.HMSet(user.ID, map[string]interface{}{
		"id":   user.ID,
		"name": user.Name,
	}).Err()
	if err != nil {
		return err
	}
	return nil
}

func getUser(client *redis.Client, id string) (User, error) {
	val, err := client.HGetAll(id).Result()
	if err != nil {
		return User{}, err
	}
	return User{
		ID:   val["id"],
		Name: val["name"],
	}, nil
}

func updateUser(client *redis.Client, user User) error {
	err := client.HMSet(user.ID, map[string]interface{}{
		"id":   user.ID,
		"name": user.Name,
	}).Err()
	if err != nil {
		return err
	}
	return nil
}

func deleteUser(client *redis.Client, id string) error {
	err := client.Del(id).Err()
	if err != nil {
		return err
	}
	return nil
}

func getAllData(client *redis.Client) error {
	all := client.Scan(0, "*", 0).Iterator()
	for all.Next() {
		key := all.Val()

		val, err := client.Get(key).Result()
		if err != nil {
			return err
		} else {
			fmt.Printf("key : %s, val %s\n", key, val)
		}
	}

	return nil
}

func main() {
	client := createClient()
	defer client.Close()

	// err := setKey(client, "user_id", "123")
	// if err != nil {
	// 	log.Fatalf("Failed to set key: %v", err)
	// }

	// val, err := getKey(client, "user_id")
	// if err != nil {
	// 	log.Fatalf("Failed to get key: %v", err)
	// }

	// fmt.Println("user_id", val)

	user := User{
		ID:   "123",
		Name: "John Doe",
	}

	err := createUser(client, user)
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	err = updateUser(client, User{
		ID:   "123",
		Name: "Jane Doe",
	})

	user, err = getUser(client, "123")
	if err != nil {
		log.Fatalf("Failed to get user: %v", err)
	}

	err = deleteUser(client, "123")

	// fmt.Println("User", user)

	err = getAllData(client)
	if err != nil {
		log.Fatalf("Failed to get all users: %v", err)
	}

	// fmt.Println("Users", users)
}
