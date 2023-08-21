package main

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

const RequestId = "request-id"

func waitingLine(w http.ResponseWriter, r *http.Request) {
	requestId := getRequestIdFromHeader(w.Header())
	if requestId == "" {
		return
	}
	addLine(requestId)
}

func checkCanEnter() bool {
	// TODO 메인 서버에 입장 가능한 사용자의 수를 전송 받아야 함
	var enterableUserNum int64 = 100

}

func getRequestIdFromHeader(h http.Header) string {
	return h.Get(RequestId)
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

	mux.HandleFunc("/favicon.ico", doNothing)
	http.ListenAndServe(":80", mux)
}
