package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"strconv"
)

const (
	RedisAddr     = "localhost:6379"
	RedisPassword = ""

	NextWaitingNumTopic = "next_waiting_num"
	EntryNumberTopic    = "entry_num"
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

func PutWaitingNum(uuid string) int64 {
	nextWaitingNum := getWaitingNum()
	client, ctx := connRedis()
	client.Append(ctx, uuid, strconv.FormatInt(nextWaitingNum, 10))

	return nextWaitingNum
}

func GetWaitingNumByRequestId(uuid string) int64 {
	client, ctx := connRedis()
	result, err := client.Get(ctx, uuid).Result()
	if err != nil {
		return PutWaitingNum(uuid)
	}
	ret, _ := strconv.ParseInt(result, 10, 64)
	return ret
}

func GetEntryNum() int64 {
	client, ctx := connRedis()
	result, err := client.Get(ctx, EntryNumberTopic).Result()
	if err != nil {
		client.IncrBy(ctx, EntryNumberTopic, 1)
	}
	ret, _ := strconv.ParseInt(result, 10, 64)
	return ret
}
