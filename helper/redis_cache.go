package helper

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
)

type Author struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func RedisClient() *redis.Client {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	json, err := json.Marshal(Author{Name: "Aamir", Age: 25})
	if err != nil {
		fmt.Println(err)
	}

	err = client.Set("id1234", json, 0).Err()
	if err != nil {
		fmt.Println(err)
	}

	val, err := client.Get("id1234").Result()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(val)
	//return rdb
	return nil
}
