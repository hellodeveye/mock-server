package server

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
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

	return http.ListenAndServe(":"+strconv.FormatUint(server.Port, 10), mux)
}

func (s *Server) createNacosClient() (naming_client.INamingClient, error) {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(s.Config.Nacos.ServerAddr, s.Config.Nacos.ServerPort),
	}

	cc := constant.ClientConfig{
		NamespaceId:         s.Config.Nacos.NamespaceId,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
	}

	return clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
}
