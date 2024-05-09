package handler

import (
	"mock-server/config"
	"net/http"
	"time"
)

func MakeHandler(ep config.Endpoint, method string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		time.Sleep(time.Duration(ep.Delay) * time.Millisecond)

		for key, value := range ep.Headers {
			w.Header().Set(key, value)
		}
		w.WriteHeader(ep.Status)
		w.Write([]byte(ep.Response))
	}
}
