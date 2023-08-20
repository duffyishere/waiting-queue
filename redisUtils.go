package main

import (
	"context"
	"github.com/redis/go-redis/v9"
)

const (
	RedisAddr     = "localhost:6379"
	RedisPassword = ""
	Topic         = "test"
)

func connRedis() (*redis.Client, context.Context) {
	client := redis.NewClient(&redis.Options{
		Addr:     RedisAddr,
		Password: RedisPassword,
		DB:       0,
	})
	ctx := context.Background()

	return client, ctx
}

func addLine(requestId string) {
	client, ctx := connRedis()
	err := client.ZAdd(ctx, Topic, redis.Z{Score: 1, Member: requestId}).Err()
	if err != nil {
		panic(err)
	}
}

func getRequestN(start, stop int64) []string {
	client, ctx := connRedis()
	result, err := client.ZRange(ctx, Topic, start, stop).Result()
	if err != nil {
		panic(err)
	}
	return result
}
