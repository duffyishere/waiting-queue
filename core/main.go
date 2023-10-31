package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	pollingHandler := http.HandlerFunc(Polling)
	mux.Handle("/p", SetContentTypeJsonMiddleware(pollingHandler))

	mux.HandleFunc("/favicon.ico", doNothing)

	http.ListenAndServe(":80", mux)
}

const RequestIdHeaderKey = "request-id"

var key = "12345678901234567890123456789012"
var iv = "1234567890123456"

type Response struct {
	RequestId string `json:"requestId"`
	Ticket    string `json:"ticket"`
}

func Polling(w http.ResponseWriter, r *http.Request) {
	uuid := GetRequestIdFromHeader(r.Header)
	client, ctx := connRedis()
	response := Response{}

	if IsAlreadyWaiting(client, ctx, uuid) {
		if !CanEnter(client, ctx, uuid) {
			response = Response{RequestId: uuid, Ticket: ""}
		} else {
			response = Response{RequestId: uuid, Ticket: ticketing(uuid)}
		}
	} else {
		AddWaitingLine(client, ctx, uuid)
		response = Response{RequestId: uuid, Ticket: ""}
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	w.Write(responseJson)
}

func ticketing(uuid string) string {
	return Ase256Encode(uuid, key, iv, 64)
}

func GetRequestIdFromHeader(h http.Header) string {
	requestId := h.Get(RequestIdHeaderKey)
	if requestId == "" {
		requestId = uuid.NewString()
		h.Set(RequestIdHeaderKey, requestId)
	}
	return requestId
}

func doNothing(w http.ResponseWriter, r *http.Request) {}

func SetContentTypeJsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
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
