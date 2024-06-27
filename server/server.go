package server

import (
	"context"
	"log"
	"mock-server/config"
	"mock-server/handler"
	"net/http"
	"strconv"
)

type MockServer struct {
	Config *config.Config
}

func New() *MockServer {
	return &MockServer{}
}

func (s *MockServer) Init(cfgPath string) {
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	err = cfg.InitConfig()
	if err != nil {
		log.Fatalf("Failed to init configuration: %v", err)
	}
	s.Config = cfg
}

func (s *MockServer) Run() {
	for _, server := range s.Config.Server {
		go doRunMockServer(*s, *server)
	}
	select {}
}

func doRunMockServer(s MockServer, service config.Server) {
	if s.Config.Nacos.Enabled {
		if err := RegisterWithNacos(s.Config.Nacos, service.Name, service.Port); err != nil {
			log.Fatalf("Failed to register with Nacos: %s", err)
		}
	}
	mux := http.NewServeMux()
	ctx := context.WithValue(context.Background(), "router", s.Config.Router)
	ctx = context.WithValue(ctx, "server", service)
	for _, ep := range service.Endpoints {
		if !ep.Enabled {
			continue
		}
		initializeRoutes(service, *ep, mux)
	}
	err := http.ListenAndServe(":"+strconv.FormatInt(int64(service.Port), 10), LoggingMiddleware(NotFoundMiddleware(mux, ctx)))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initializeRoutes(server config.Server, ep config.Endpoint, mux *http.ServeMux) {
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
