package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type RedisClient struct {
	c *redis.Client
}

var (
	client                 = &RedisClient{}
	ctx    context.Context = context.Background()
)

func init() {
	//Making the connection to redis server
	c := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	if err := c.Ping(ctx).Err(); err != nil {
		panic("Error connecting to database: Reason: " + err.Error())
	}
	client.c = c
}
func SetKey(key string, value interface{}, expiration time.Duration) error {
	setErr := client.c.Set(ctx, key, value, expiration).Err()
	if setErr != nil {
		log.Fatal("Error setting the keys: Reason ", setErr)
	}
	fmt.Println("Successfully set the keys")
	return nil
}
func GetKey(key string) error {
	var (
		getError error
		value    interface{}
	)
	value, getError = client.c.Get(ctx, key).Result()
	if getError != nil {
		log.Fatal("Error getting the value for the specified key: Reason: ", getError)
	}
	fmt.Println("Successfully retrieved the keys")
	fmt.Println("Value is : ", value)
	return nil
}
