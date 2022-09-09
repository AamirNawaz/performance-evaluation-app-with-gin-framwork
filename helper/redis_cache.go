package helper

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var ctx = context.Background()

func RedisClient() *redis.Client {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return client
}

func SetExVal(key string, val string, exp time.Duration) error {
	_, err := RedisClient().Set(context.Background(), key, val, exp).Result()
	if err != nil {
		fmt.Println("inside setVal func", err)
	}
	return nil
}

func GetExVal(key string) string {
	value, err := RedisClient().Get(context.Background(), key).Result()
	if err != nil {
		fmt.Println("inside setVal func", err)
	}
	return value
}
