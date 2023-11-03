package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	WaitingLineTopic = "waiting_line"
	RunningMapTopic  = "running_map"
	EntryNumberTopic = "entry_num"
)

func ConnRedis() (*redis.Client, context.Context) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPassword,
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
	num = num - 1
	values, err := client.ZRange(ctx, WaitingLineTopic, 0, num).Result()
	if err != nil {
		panic(err)
	}
	removeForWaitingLine(client, ctx, num)
	addRunningMap(client, ctx, values)
}

func addRunningMap(client *redis.Client, ctx context.Context, uuids []string) {
	for i := range uuids {
		println("UUID: ", uuids[i])
		err := client.HSet(ctx, RunningMapTopic, uuids[i], 0).Err()
		if err != nil {
			panic(err)
		}
	}
}

func removeForWaitingLine(client *redis.Client, ctx context.Context, num int64) {
	err := client.ZRemRangeByRank(ctx, WaitingLineTopic, 0, num).Err()
	if err != nil {
		panic(err)
	}
}
