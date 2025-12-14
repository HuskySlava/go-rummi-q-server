package main

import (
	"go-rummi-q-server/internal/config"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed loading config", err)
	}
	log.Println("CFG:", cfg)
}
