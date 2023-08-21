package main

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

const RequestId = "request-id"

// TODO: 별도의 쓰레드로 항상 Redis 와 동기화 되어있어야 함
var userCapacity int64

func waitingLine(w http.ResponseWriter, r *http.Request) {
	requestId := getRequestIdFromHeader(w.Header())
	addLine(requestId)
}

func polling(w http.ResponseWriter, r *http.Request) {
	requestId := getRequestIdFromHeader(w.Header())
	if checkCanEnter(requestId) {
		// TODO: 메인 서버의 입장권을 발급 후 전송
	}
	else {
		// TODO: 입장 실패 메시지 전송
	}
}

func checkCanEnter(requestId string) bool {
	waitingNum := getWaitingNumBy(requestId)
	return waitingNum < userCapacity
}

func getRequestIdFromHeader(h http.Header) string {
	requestId := h.Get(RequestId)
	if requestId == "" {
		panic("해당 요청의 request-id가 존재하지 않습니다.")
	}
	return requestId
}

func doNothing(w http.ResponseWriter, r *http.Request) {}

func generateRequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		xRequestID := uuid.New().String()
		w.Header().Set(RequestId, xRequestID)
		next.ServeHTTP(w, r)
	})
}

func logMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf(">>> Request: %s %s %s\n", r.Host, r.URL.Path, w.Header().Get("request-id"))
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()

	finalHandler := http.HandlerFunc(waitingLine)
	mux.Handle("/", generateRequestId(logMiddleWare(finalHandler)))

	mux.HandleFunc("/polling", polling)
	mux.HandleFunc("/favicon.ico", doNothing)
	http.ListenAndServe(":80", mux)
}
