package transport

import (
	"log"
	"net/http"
)

func StartServer() {
	mux := NewRouter()

	if err := http.ListenAndServe(":3001", mux); err != nil {
		log.Printf("Server failed: %v", err)
	}
}
