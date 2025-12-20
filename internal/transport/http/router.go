package transport

import (
	"encoding/json"
	"fmt"
	"go-rummi-q-server/internal/domain/lobby"
	"net/http"
	"strings"
)

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		status, err := fmt.Fprint(w, "Welcome to Rummi-Q-Server")
		fmt.Println("Status:", status, "Error", err)
	})

	mux.HandleFunc("POST /lobby", func(w http.ResponseWriter, r *http.Request) {
		gameLobby := lobby.NewLobby()
		fmt.Println(gameLobby)

		resp := map[string]any{
			"game_id": gameLobby.ID,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	mux.HandleFunc("POST /lobby/", func(w http.ResponseWriter, r *http.Request) {

		path := r.URL.Path
		parts := strings.Split(path, "/")
		// Validate path structure "/games/{id}/join"
		if len(parts) != 4 || parts[3] != "join" {
			http.NotFound(w, r)
			return
		}

		if r.Body == nil {
			http.Error(w, "Request body required", http.StatusBadRequest)
			return
		}

		var req struct {
			PlayerName string `json:"player_name"`
		}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		if req.PlayerName == "" {
			http.Error(w, "`player_name` is required", http.StatusBadRequest)
			return
		}

	})

	return mux
}
