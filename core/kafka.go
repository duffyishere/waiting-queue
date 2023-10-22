package main

import "github.com/confluentinc/confluent-kafka-go/v2/kafka"

const (
	RequestIdHeaderKey = "request-id"

	KafkaAddr    = "localhost:29092"
	KafkaGroupId = "myGroup"
)

var KafkaTopicNames = []string{
	"streaming.extra-user-capacity-num",
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
