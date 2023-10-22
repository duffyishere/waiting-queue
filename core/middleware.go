package main

import (
	"crypto/aes"
	"fmt"
	"net/http"
	"strconv"
)

func doNothing(w http.ResponseWriter, r *http.Request) {}

func SetContentTypeJsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func GenerateRequestIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := Ase256Encode(strconv.FormatInt(GetWaitingNum(), 10), key, iv, aes.BlockSize)
		w.Header().Set(RequestIdHeaderKey, requestId)
		next.ServeHTTP(w, r)
	})
}

func LogMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf(">>> Host: %s, URL: %s, RequestID: %s\n", r.Host, r.URL.Path, GetRequestIdFromHeader(w.Header()))
		next.ServeHTTP(w, r)
	})
}
