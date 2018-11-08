package main

import (
	"fmt"

	"github.com/go-redis/redis"
)

func keys(redisdb *redis.Client) {

	keys, err := redisdb.Keys("*alu*").Result()
	if err != nil {
		fmt.Println("keys: ", err)
		return
	}

	fmt.Println("keys:", keys)
}

func main() {
	redisdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := redisdb.Ping().Result()
	fmt.Println(pong, err)

	redisdb.FlushDB()
	for i := 0; i < 33; i++ {
		if i%2 == 1 {
			err := redisdb.Set(fmt.Sprintf("alu%d", i), "value", 0).Err()
			if err != nil {
				panic(err)
			}

			continue
		}

		err := redisdb.Set(fmt.Sprintf("balu%d", i), "value", 0).Err()
		if err != nil {
			panic(err)
		}
	}

	keys(redisdb)

	var cursor uint64
	var n int
	for {
		var keys []string
		var err error
		keys, cursor, err = redisdb.Scan(cursor, "alu*", 10).Result()

		if err != nil {
			panic(err)
		}
		n += len(keys)

		fmt.Println("keys:", keys, " cursor:", cursor, " err:", err, " n:", n)

		if cursor == 0 {
			break
		}
	}

	fmt.Printf("found %d keys\n", n)

}
