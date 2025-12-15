package main

import (
	"go-rummi-q-server/internal/config"
	transport "go-rummi-q-server/internal/transport/http"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed loading config:, %v", err)
	}
	log.Println("CFG:", cfg)

	transport.StartServer()
}
