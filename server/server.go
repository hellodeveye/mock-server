package server

import (
	"mock-server/config"
	"mock-server/handler"
	"net/http"
	"strconv"
)

type Server struct {
	Config *config.Config
}

func New(cfg *config.Config) *Server {
	return &Server{Config: cfg}
}

func (s *Server) Run(server config.Server) error {

	mux := http.NewServeMux()
	for _, ep := range server.Endpoints {
		if !ep.Enabled {
			continue
		}
		switch ep.Method {
		case "GET":
			mux.HandleFunc(ep.URL, handler.MakeHandler(ep, http.MethodGet))
		case "POST":
			mux.HandleFunc(ep.URL, handler.MakeHandler(ep, http.MethodPost))
		case "PUT":
			mux.HandleFunc(ep.URL, handler.MakeHandler(ep, http.MethodPut))
		case "DELETE":
			mux.HandleFunc(ep.URL, handler.MakeHandler(ep, http.MethodDelete))
		}
	}

	return http.ListenAndServe(":"+strconv.FormatUint(server.Port, 10), LoggingMiddleware(mux))
}
