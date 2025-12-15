package main

import (
	"go-rummi-q-server/internal/config"
	transport "go-rummi-q-server/internal/transport/http"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed loading config: %v", err)
	}

	log.Printf("Starting server at %s:%d", cfg.HTTPConfig.ListenHost, cfg.HTTPConfig.ListenPort)
	err = transport.StartServer(cfg.HTTPConfig)
	if err != nil {
		log.Fatalf("Server start error: %v", err)
	}

}
