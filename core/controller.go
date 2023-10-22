package main

import "net/http"

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
		panic("The request-id for that request does not exist.")
	}
	return requestId
}
