package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	RedisAddr        = "localhost:6379"
	RedisPassword    = ""
	TopicExpiredTime = time.Hour * 2

	NextWaitingNumTopic = "next_waiting_num"
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

func getWaitingNum() int64 {
	client, ctx := connRedis()
	result, err := client.IncrBy(ctx, NextWaitingNumTopic, 1).Result()
	if err != nil {
		panic(err)
	}
	return result
}
