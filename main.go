package main

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"time"
)

const (
	RequestIdHeaderKey = "request-id"

	KafkaAddr    = "localhost:29092"
	KafkaGroupId = "myGroup"
)

var KafkaTopicNames = []string{
	"streaming.extra-user-capacity-num",
}

var userCapacity int64

func waitingQueue(w http.ResponseWriter, r *http.Request) {
	requestId := getRequestIdFromHeader(w.Header())
	addQueue(requestId)
}

type PollingResponse struct {
	Result        bool   `json:"result"`
	Ticket        string `json:"ticket"`
	NowWaitingNum int64  `json:"now_waiting_num"`
}

func polling(w http.ResponseWriter, r *http.Request) {
	requestId := getRequestIdFromHeader(r.Header)
	waitingNum := getWaitingNumBy(requestId)
	var data PollingResponse
	if waitingNum < userCapacity {
		// TODO: 암호화 된 입장권 생성
		enterTicket := "test1234"
		data = PollingResponse{
			Result:        true,
			Ticket:        enterTicket,
			NowWaitingNum: 0,
		}
	} else {
		data = PollingResponse{
			Result:        false,
			Ticket:        "",
			NowWaitingNum: waitingNum,
		}
	}
	jsonData, err := json.Marshal(data)
	fmt.Println(string(jsonData))
	if err != nil {
		panic(err)
	}

	w.Write(jsonData)
}

func getRequestIdFromHeader(h http.Header) string {
	requestId := h.Get(RequestIdHeaderKey)
	if requestId == "" {
		panic("The request-id for that request does not exist.")
	}
	return requestId
}

func doNothing(w http.ResponseWriter, r *http.Request) {}

func setContentTypeJsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func generateRequestIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := uuid.New().String()
		w.Header().Set(RequestIdHeaderKey, requestId)
		next.ServeHTTP(w, r)
	})
}

func logMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf(">>> Host: %s, URL: %s, RequestID: %s\n", r.Host, r.URL.Path, getRequestIdFromHeader(w.Header()))
		next.ServeHTTP(w, r)
	})
}

func main() {
	go updateUserCapacity(connectKafka())

	mux := http.NewServeMux()

	finalHandler := http.HandlerFunc(waitingQueue)
	mux.Handle("/", generateRequestIdMiddleware(logMiddleWare(finalHandler)))

	pollingHandler := http.HandlerFunc(polling)
	mux.Handle("/polling", setContentTypeJsonMiddleware(pollingHandler))

	mux.HandleFunc("/favicon.ico", doNothing)

	http.ListenAndServe(":80", mux)
}

func connectKafka() *kafka.Consumer {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": KafkaAddr,
		"group.id":          KafkaGroupId,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}
	return consumer
}

func updateUserCapacity(consumer *kafka.Consumer) {
	consumer.SubscribeTopics(KafkaTopicNames, nil)
	for {
		msg, err := consumer.ReadMessage(time.Second)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			additionalUserCapacity, _ := strconv.ParseInt(string(msg.Value), 10, 64)
			userCapacity = increaseUserCapacity(additionalUserCapacity)
			fmt.Printf("Now user capacity: %d\n", userCapacity)
		} else if !err.(kafka.Error).IsTimeout() {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}
