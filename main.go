package main

import (
	"context"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("gEdbih-reqfyf-1p1")
	store = sessions.NewCookieStore(key)
)

const (
	RedisAddr     = "localhost:6379"
	RedisPassword = ""
	ExpiredTime   = time.Hour
)

func myHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session_id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if session.Values["ticket"] == nil {
		nextTicketNum := increaseTicketCount(connRedis())
		// TODO: 티켓의 만료 기간을 설정해야 함
		session.Values["ticket"] = nextTicketNum
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func connRedis() (*redis.Client, context.Context) {
	client := redis.NewClient(&redis.Options{
		Addr:     RedisAddr,
		Password: RedisPassword,
		DB:       0,
	})
	ctx := context.Background()

	return client, ctx
}

func increaseTicketCount(client *redis.Client, ctx context.Context) int64 {
	count, err := client.Incr(ctx, "ticket_count").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(count)
	return count
}

func main() {
	http.HandleFunc("/", myHandler)
	http.ListenAndServe(":80", nil)
}
