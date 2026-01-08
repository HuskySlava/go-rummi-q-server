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
	writeResponse(w, http.StatusCreated, map[string]any{
		"message": "Lobby Created",
		"game_id": gameLobby.ID,
	})
}

func joinLobby(w http.ResponseWriter, r *http.Request, lobbyID uuid.UUID) {
	var req struct {
		PlayerName string `json:"player_name"`
		PlayerID   string `json:"player_id"`
	}

	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}
	if req.PlayerName == "" {
		writeError(w, http.StatusBadRequest, "`player_name` is required")
		return
	}
	if !lobbies.LobbyExists(lobbyID) {
		writeError(w, http.StatusNotFound, "Lobby not found")
		return
	}
	if err := lobbies.JoinLobby(lobbyID, req.PlayerID, req.PlayerName); err != nil {
		writeError(w, http.StatusBadRequest, "Unable to join lobby: "+err.Error())
		return
	}

	writeResponse(w, http.StatusOK, map[string]any{
		"message":     "Player joined game",
		"player_name": req.PlayerName,
		"lobby_id":    lobbyID.String(),
	})
}

func newPlayer(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PlayerName string `json:"player_name"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}
	if req.PlayerName == "" {
		writeError(w, http.StatusBadRequest, "`player_name` is required")
		return
	}

	id, err := game.GeneratePlayerID()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to generate player ID: "+err.Error())
		return
	}
	game.AddPlayer(id, req.PlayerName)

	writeResponse(w, http.StatusCreated, map[string]any{
		"message":     "New player created",
		"player_name": req.PlayerName,
		"player_id":   string(id[:]),
	})
}

// getAllPlayers is a temporary HTTP handler for debugging.
// TODO: restrict or remove before release.
func getAllPlayers(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	urlParts := strings.Split(path, "/")

	w.Header().Set("Content-Type", "application/json")

	// Expecting: /players
	if len(urlParts) != 1 || urlParts[0] != "players" {
		http.NotFound(w, r)
		return
	}
	resp, err := game.GetAllPlayersJSON()
	if err != nil {
		http.Error(w, "Failed to get player list: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}

// ## Router handlers ##

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /lobbies", createLobby)
	mux.HandleFunc("POST /player", newPlayer)
	mux.HandleFunc("GET /players", getAllPlayers)
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

	mux.HandleFunc("GET /", defaultRoute)

	return mux
}

func writeResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeResponse(w, status, map[string]string{"error": message})
}

func decodeJSON(r *http.Request, dest any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(dest)
}
