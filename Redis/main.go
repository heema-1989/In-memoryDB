package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"strings"
	"time"
)

func main() {
	//Making the connection to the redis server
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", //no password set
		DB:       0,  //use default db-->this could be anything 1,2,3,4
	})
	ctx := context.Background()
	fmt.Println(strings.Repeat("-", 20))
	//Setting the key-valued pair
	err := client.Set(ctx, "username", "Heema", 0).Err()
	if err != nil {
		log.Fatal("Error setting keys: ", err)
	}
	//Retrieving the key-valued pair--here if we write a key name that does not exist, then it will throw runtime error
	val, getErr := client.Get(ctx, "username").Result()
	if getErr != nil {
		log.Fatal("error getting the key-values: ", err)
	}
	fmt.Println("Value1:", val)
	fmt.Println(strings.Repeat("-", 20))
	val2, err := client.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	fmt.Println(strings.Repeat("-", 20))
	err = client.Set(ctx, "username", "Dhatri", 0).Err()
	if err != nil {
		log.Fatal("Error setting key-values: ", err)
	}
	//--Here if we give another value with same key name then it will overwrite the previous
	//val3, err := client.Get(ctx, "username").Result()
	//if err != nil {
	//	log.Fatal("Err getting key-values: ", err)
	//}
	//fmt.Println("Value3: ", val3)
	fmt.Println(strings.Repeat("-", 20))
	set1, err := client.SetNX(ctx, "key", "value", 10*time.Second).Result()
	fmt.Println(set1)
	fmt.Println(strings.Repeat("-", 20))
	set2, err := client.Set(ctx, "key2", "value2", redis.KeepTTL).Result()
	fmt.Println(set2)
	//Checking whether the connection is successfully established or not--equivalent to redis-cli ping for normal redis-cli
	pong, err := client.Ping(ctx).Result()
	fmt.Println(pong, err)

}
