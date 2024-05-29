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
	err = cfg.InitConfig()
	if err != nil {
		log.Fatalf("Failed to init configuration: %v", err)
	}
	enabled := cfg.Nacos.Enabled
	for _, service := range cfg.Server {
		go func(service config.Server) {
			if enabled {
				if err := server.RegisterWithNacos(cfg.Nacos, service.Name, service.Port); err != nil {
					log.Fatalf("Failed to register with Nacos: %s", err)
				}
			}
			srv := server.New(cfg)
			if err := srv.Run(service); err != nil {
				log.Fatalf("Failed to run server: %v", err)
			}
		}(service)
	}
	select {}

}
