package server

import (
	"log"
	"net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("请求方法: %s, 请求URL: %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
