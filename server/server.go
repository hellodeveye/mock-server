package server

import (
	"context"
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
	ctx := context.WithValue(context.Background(), "router", s.Config.Router)
	ctx = context.WithValue(ctx, "server", server)
	for _, ep := range server.Endpoints {
		if !ep.Enabled {
			continue
		}
		switch ep.Method {
		case "GET":
			mux.HandleFunc(ep.URL, handler.MakeHandler(ep, server, http.MethodGet))
		case "POST":
			mux.HandleFunc(ep.URL, handler.MakeHandler(ep, server, http.MethodPost))
		case "PUT":
			mux.HandleFunc(ep.URL, handler.MakeHandler(ep, server, http.MethodPut))
		case "DELETE":
			mux.HandleFunc(ep.URL, handler.MakeHandler(ep, server, http.MethodDelete))
		}
	}
	return http.ListenAndServe(":"+strconv.FormatInt(int64(server.Port), 10), LoggingMiddleware(NotFoundMiddleware(mux, ctx)))
}
