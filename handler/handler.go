package handler

import (
	"context"
	"mock-server/config"
	"net/http"
	"time"
)

func MakeHandler(ep config.Endpoint, server config.Server, method string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.WithContext(context.WithValue(context.Background(), "server", server))
		if r.Method != method {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		time.Sleep(time.Duration(ep.Delay) * time.Millisecond)

		for key, value := range ep.Headers {
			w.Header().Set(key, value)
		}
		if ep.Status == 0 {
			ep.Status = http.StatusOK
		}
		w.WriteHeader(ep.Status)
		_, _ = w.Write([]byte(ep.Response))
	}
}
