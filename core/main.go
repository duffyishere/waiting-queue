package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

var (
	requestIdHeaderKey = ""
	key                = ""
	iv                 = ""
	redisHost          = ""
	redisPassword      = ""
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Cannot find '.env' file. Please check you've added the env file.")
	}
	requestIdHeaderKey = os.Getenv("REQUEST_ID_HEADER_KEY_NAME")
	key = os.Getenv("ASE256_KEY")
	iv = os.Getenv("ASE256_IV")
	redisHost = os.Getenv("REDIS_HOST")
	redisPassword = os.Getenv("REDIS_PASSWORD")

	mux := http.NewServeMux()

	pollingHandler := http.HandlerFunc(polling)
	mux.Handle("/p", setContentTypeJsonMiddleware(pollingHandler))

	mux.HandleFunc("/favicon.ico", doNothing)

	http.ListenAndServe(":3000", mux)
}

type Response struct {
	RequestId string `json:"request-id"`
	Ticket    string `json:"ticket"`
}

func polling(w http.ResponseWriter, r *http.Request) {
	uuid := getRequestIdFromHeader(r.Header)
	client, ctx := ConnRedis()
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
	return ase256Encode(uuid, key, iv, 64)
}

func getRequestIdFromHeader(h http.Header) string {
	requestId := h.Get(requestIdHeaderKey)
	if requestId == "" {
		requestId = uuid.NewString()
		h.Set(requestIdHeaderKey, requestId)
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

func ase256Encode(plaintext string, key string, iv string, blockSize int) string {
	bKey := []byte(key)
	bIV := []byte(iv)
	bPlaintext := PKCS5Padding([]byte(plaintext), blockSize, len(plaintext))
	block, _ := aes.NewCipher(bKey)
	ciphertext := make([]byte, len(bPlaintext))
	mode := cipher.NewCBCEncrypter(block, bIV)
	mode.CryptBlocks(ciphertext, bPlaintext)
	return hex.EncodeToString(ciphertext)
}

func ase256Decode(cipherText string, encKey string, iv string) (decryptedString string) {
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
