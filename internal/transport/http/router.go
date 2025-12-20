package transport

import (
	"encoding/json"
	"fmt"
	"go-rummi-q-server/internal/domain/lobby"
	"net/http"
)

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		status, err := fmt.Fprint(w, "Welcome to Rummi-Q-Server")
		fmt.Println("Status:", status, "Error", err)
	})

	mux.HandleFunc("Post /games", func(w http.ResponseWriter, r *http.Request) {
		gameLobby := lobby.NewLobby()
		fmt.Println(gameLobby)

		resp := map[string]any{
			"game_id": gameLobby.ID,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	mux.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	})

	return mux
}
