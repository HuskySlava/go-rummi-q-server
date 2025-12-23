package transport

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go-rummi-q-server/internal/domain/lobbies"
	"net/http"
	"strings"
)

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// Test route
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		status, err := fmt.Fprint(w, "Welcome to Rummi-Q-Server")
		fmt.Println("Status:", status, "Error", err)
	})

	// Create new lobby
	mux.HandleFunc("POST /lobbies", func(w http.ResponseWriter, r *http.Request) {
		gameLobby := lobbies.NewLobby()

		resp := map[string]any{
			"message": "Lobby Created",
			"game_id": gameLobby.ID,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	// Join existing lobby
	mux.HandleFunc("POST /lobbies/", func(w http.ResponseWriter, r *http.Request) {

		path := r.URL.Path
		parts := strings.Split(path, "/")
		// Validate path structure "/games/{id}/join"
		if len(parts) != 4 || parts[3] != "join" {
			http.NotFound(w, r)
			return
		}
		lobbyId, err := uuid.Parse(parts[2])
		if err != nil {
			http.Error(w, "Invalid lobby ID", http.StatusBadRequest)
			return
		}

		if r.Body == nil {
			http.Error(w, "Request body required", http.StatusBadRequest)
			return
		}

		var req struct {
			PlayerName string `json:"player_name"`
		}

		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		if req.PlayerName == "" {
			http.Error(w, "`player_name` is required", http.StatusBadRequest)
			return
		}

		// Join logic
		if !lobbies.LobbyExists(lobbyId) {
			http.Error(w, "Lobby not found", http.StatusNotFound)
			return
		}

		// Success
		resp := map[string]string{
			"message":     "Player joined game",
			"player_name": req.PlayerName,
			"lobby_id":    lobbyId.String(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	return mux
}
