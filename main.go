package main

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

func waitingLine(w http.ResponseWriter, r *http.Request) {
	requestId := w.Header().Get("request-id")
	if requestId == "" {
		return
	}
	addLine(requestId)
}

func doNothing(w http.ResponseWriter, r *http.Request) {}

func generateRequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		xRequestID := uuid.New().String()
		w.Header().Set("request-id", xRequestID)
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
