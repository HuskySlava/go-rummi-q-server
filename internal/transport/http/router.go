package transport

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go-rummi-q-server/internal/domain/game"
	"go-rummi-q-server/internal/domain/lobbies"
	"log"
	"net/http"
	"strings"
)

// ## Route actions ##

func defaultRoute(w http.ResponseWriter, _ *http.Request) {
	status, err := fmt.Fprint(w, "Welcome to Rummi-Q-Server")
	log.Printf("Status: %d, Error: %v \n", status, err)
}

func createLobby(w http.ResponseWriter, _ *http.Request) {
	gameLobby := lobbies.NewLobby()

	resp := map[string]any{
		"message": "Lobby Created",
		"game_id": gameLobby.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
}

func joinLobby(w http.ResponseWriter, r *http.Request, lobbyID uuid.UUID) {
	defer r.Body.Close()

	var req struct {
		PlayerName string `json:"player_name"`
		PlayerID   string `json:"player_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.PlayerName == "" {
		http.Error(w, "`player_name` is required", http.StatusBadRequest)
		return
	}

	if !lobbies.LobbyExists(lobbyID) {
		http.Error(w, "Lobby not found", http.StatusNotFound)
		return
	}

	err := lobbies.JoinLobby(lobbyID, req.PlayerID, req.PlayerName)
	if err != nil {
		http.Error(w, "Unable to join lobby", http.StatusBadRequest)
	}

	resp := map[string]string{
		"message":     "Player joined game",
		"player_name": req.PlayerName,
		"lobby_id":    lobbyID.String(),
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
}

func newPlayer(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req struct {
		PlayerName string `json:"player_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.PlayerName == "" {
		http.Error(w, "`player_name` is required", http.StatusBadRequest)
		return
	}

	id, err := game.GeneratePlayerID()
	if err != nil {
		http.Error(w, "Failed to generate player ID: "+err.Error(), http.StatusInternalServerError)
		return
	}

	game.AddPlayer(id, req.PlayerName)

	resp := map[string]string{
		"message":     "New player created",
		"player_name": req.PlayerName,
		"player_id":   string(id[:]),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
}

// ## Router handlers ##

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", defaultRoute)
	mux.HandleFunc("POST /lobbies", createLobby)
	// Lobby actions
	mux.HandleFunc("POST /lobbies/", func(w http.ResponseWriter, r *http.Request) {

		path := strings.Trim(r.URL.Path, "/")
		urlParts := strings.Split(path, "/")

		// Expecting: /lobbies/{id}/join
		if len(urlParts) != 3 || urlParts[0] != "lobbies" {
			http.NotFound(w, r)
			return
		}

		lobbyID, err := uuid.Parse(urlParts[1])
		if err != nil {
			http.Error(w, "Invalid lobby ID", http.StatusBadRequest)
			return
		}

		action := urlParts[2]
		switch action {
		case "join":
			joinLobby(w, r, lobbyID)
		default:
			http.NotFound(w, r)
		}

	})

	return mux
}
