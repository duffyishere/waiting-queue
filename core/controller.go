package main

import (
	"crypto/aes"
	"net/http"
	"strconv"
)

type PollingResponse struct {
	Result        bool   `json:"result"`
	Ticket        string `json:"ticket"`
	NowWaitingNum int64  `json:"now_waiting_num"`
}

func Polling(w http.ResponseWriter, r *http.Request) {
	requestId := GetRequestIdFromHeader(r.Header)
	println("Encrypt request id: ", requestId)
	println("Decrypt request id: ", Ase256Decode(requestId, key, iv))
}

func GetRequestIdFromHeader(h http.Header) string {
	requestId := h.Get(RequestIdHeaderKey)
	if requestId == "" {
		requestId = Ase256Encode(strconv.FormatInt(GetWaitingNum(), 10), key, iv, aes.BlockSize)
		h.Set(RequestIdHeaderKey, requestId)
	}
	return requestId
}
