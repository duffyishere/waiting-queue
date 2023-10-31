package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

const (
	RedisAddr     = "localhost:6379"
	RedisPassword = ""

	WaitingLineTopic = "waiting_line"
	RunningMapTopic  = "running_map"
	EntryNumberTopic = "entry_num"
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

func AddWaitingLine(client *redis.Client, ctx context.Context, uuid string) {
	score := float64(time.Now().Unix())
	client.ZAdd(ctx, WaitingLineTopic, redis.Z{Score: score, Member: uuid})
}

func IsAlreadyWaiting(client *redis.Client, ctx context.Context, uuid string) bool {
	err := client.ZScore(ctx, WaitingLineTopic, uuid).Err()
	if err != nil {
		return false
	}
	return true
}

func CanEnter(client *redis.Client, ctx context.Context, uuid string) bool {
	err := client.HMGet(ctx, RunningMapTopic, uuid).Err()
	if err != nil {
		return false
	}
	return true
}

func AddEntryNumber(client *redis.Client, ctx context.Context, num int64) {
	client.IncrBy(ctx, EntryNumberTopic, num)
	values, err := client.ZRange(ctx, WaitingLineTopic, 0, num).Result()
	if err != nil {
		panic(err)
	}
	removeForWaitingLine(client, ctx, num-1)
	addRunningMap(client, ctx, values)
}

func addRunningMap(client *redis.Client, ctx context.Context, uuid []string) {
	result, err := client.HMSet(ctx, RunningMapTopic, uuid).Result()
	if !result || err != nil {
		panic(err)
	}
}

func removeForWaitingLine(client *redis.Client, ctx context.Context, num int64) {
	err := client.ZRemRangeByRank(ctx, WaitingLineTopic, 0, num).Err()
	if err != nil {
		panic(err)
	}
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
