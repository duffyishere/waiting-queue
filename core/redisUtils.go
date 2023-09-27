package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

const (
	RedisAddr        = "localhost:6379"
	RedisPassword    = ""
	TopicExpiredTime = time.Hour * 2

	UserCapacity        = "user_capacity"
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

func addQueue(requestId string) {
	client, ctx := connRedis()
	nextWaitingNum := increaseWaitingNum()
	err := client.Set(ctx, requestId, nextWaitingNum, TopicExpiredTime).Err()
	if err != nil {
		panic(err)
	}
}

func getWaitingNumBy(requestId string) int64 {
	client, ctx := connRedis()
	result, err := client.Get(ctx, requestId).Result()
	if err != nil {
		panic(err)
	}
	waitingNum, err := strconv.ParseInt(result, 10, 64)
	if err != nil {
		panic(err)
	}
	return waitingNum
}

func getUserCapacity() int64 {
	client, ctx := connRedis()
	result, err := client.IncrBy(ctx, UserCapacity, 0).Result()
	if err != nil {
		panic(err)
	}
	return result
}

func increaseUserCapacity(value int64) int64 {
	client, ctx := connRedis()
	result, err := client.IncrBy(ctx, UserCapacity, value).Result()
	if err != nil {
		panic(err)
	}
	return result
}

func increaseWaitingNum() int64 {
	client, ctx := connRedis()
	result, err := client.IncrBy(ctx, NextWaitingNumTopic, 1).Result()
	if err != nil {
		panic(err)
	}
	return result
}
