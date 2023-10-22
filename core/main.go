package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"net/http"
	"strconv"
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
var key = "12345678901234567890123456789012"
var iv = "1234567890123456"

type PollingResponse struct {
	Result        bool   `json:"result"`
	Ticket        string `json:"ticket"`
	NowWaitingNum int64  `json:"now_waiting_num"`
}

func polling(w http.ResponseWriter, r *http.Request) {
	requestId := getRequestIdFromHeader(r.Header)
	println("Encrypt request id: ", requestId)
	println("Decrypt request id: ", Ase256Decode(requestId, key, iv))
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
		requestId := Ase256Encode(strconv.FormatInt(getWaitingNum(), 10), key, iv, aes.BlockSize)
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
	mux := http.NewServeMux()
	finalHandler := http.HandlerFunc(doNothing)
	mux.Handle("/", generateRequestIdMiddleware(logMiddleWare(finalHandler)))

	pollingHandler := http.HandlerFunc(polling)
	mux.Handle("/p", setContentTypeJsonMiddleware(pollingHandler))

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

func Ase256Encode(plaintext string, key string, iv string, blockSize int) string {
	bKey := []byte(key)
	bIV := []byte(iv)
	bPlaintext := PKCS5Padding([]byte(plaintext), blockSize, len(plaintext))
	block, _ := aes.NewCipher(bKey)
	ciphertext := make([]byte, len(bPlaintext))
	mode := cipher.NewCBCEncrypter(block, bIV)
	mode.CryptBlocks(ciphertext, bPlaintext)
	return hex.EncodeToString(ciphertext)
}

func Ase256Decode(cipherText string, encKey string, iv string) (decryptedString string) {
	bKey := []byte(encKey)
	bIV := []byte(iv)
	cipherTextDecoded, err := hex.DecodeString(cipherText)
	if err != nil {
		panic(err)
	}

	block, err := aes.NewCipher(bKey)
	if err != nil {
		panic(err)
	}

	mode := cipher.NewCBCDecrypter(block, bIV)
	mode.CryptBlocks([]byte(cipherTextDecoded), []byte(cipherTextDecoded))
	return string(cipherTextDecoded)
}

func PKCS5Padding(ciphertext []byte, blockSize int, after int) []byte {
	padding := (blockSize - len(ciphertext)%blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
