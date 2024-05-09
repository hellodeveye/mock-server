package main

import (
	"flag"
	"log"
	"mock-server/config"
	"mock-server/server"
)

func main() {
	configPath := flag.String("config", "config.yml", "Path to the configuration file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	if cfg.Nacos.Enabled {

		for _, service := range cfg.Server {
			go func(service config.Server) {
				if err := server.RegisterWithNacos(cfg.Nacos, service.Name, service.Port); err != nil {
					log.Fatalf("Failed to register with Nacos: %s", err)
				}
			}(service)
		}
	}

	for _, service := range cfg.Server {
		go func(service config.Server) {
			srv := server.New(cfg)
			if err := srv.Run(service); err != nil {
				log.Fatalf("Failed to run server: %v", err)
			}
		}(service)
	}
	select {}

}
