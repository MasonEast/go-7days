package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

func TestRedis(t *testing.T) {
	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
			t.Fatal(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
			t.Fatal(err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
			fmt.Println("key2 does not exist")
	} else if err != nil {
			t.Fatal(err)
	} else {
			fmt.Println("key2", val2)
	}
}


