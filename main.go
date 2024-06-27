package main

import (
	"flag"
	"mock-server/server"
)

func main() {
	configPath := flag.String("config", "config.yml", "Path to the configuration file")
	flag.Parse()

	srv := server.New()
	srv.Init(*configPath)
	srv.Run()
}
