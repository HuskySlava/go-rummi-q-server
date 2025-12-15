package transport

import (
	"fmt"
	"net/http"
)

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		status, err := fmt.Fprint(w, "Welcome to Rummi-Q-Server")
		fmt.Println("Status:", status, "Error", err)
	})

	return mux
}
