package transport

import (
	"fmt"
	"github.com/google/uuid"
	"go-rummi-q-server/internal/domain/game"
	"go-rummi-q-server/internal/domain/lobbies"
	"log"
	"net/http"
	"strings"
)

func defaultRoute(w http.ResponseWriter, _ *http.Request) {
	status, err := fmt.Fprint(w, "Welcome to Rummi-Q-Server")
	log.Printf("Status: %d, Error: %v \n", status, err)
}

func createLobby(w http.ResponseWriter, _ *http.Request) {
	gameLobby := lobbies.NewLobby()
	WriteResponse(w, http.StatusCreated, map[string]any{
		"message": "Lobby Created",
		"game_id": gameLobby.ID,
	})
}

func joinLobby(w http.ResponseWriter, r *http.Request, lobbyID uuid.UUID) {
	var req struct {
		PlayerName string `json:"player_name"`
		PlayerID   string `json:"player_id"`
	}

	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}
	if req.PlayerName == "" {
		WriteError(w, http.StatusBadRequest, "`player_name` is required")
		return
	}
	if !lobbies.LobbyExists(lobbyID) {
		WriteError(w, http.StatusNotFound, "Lobby not found")
		return
	}
	if err := lobbies.JoinLobby(lobbyID, req.PlayerID, req.PlayerName); err != nil {
		WriteError(w, http.StatusBadRequest, "Unable to join lobby: "+err.Error())
		return
	}

	WriteResponse(w, http.StatusOK, map[string]any{
		"message":     "Player joined game",
		"player_name": req.PlayerName,
		"lobby_id":    lobbyID.String(),
	})
}

func newPlayer(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PlayerName string `json:"player_name"`
	}
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}
	if req.PlayerName == "" {
		WriteError(w, http.StatusBadRequest, "`player_name` is required")
		return
	}

	id, err := game.GeneratePlayerID()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to generate player ID: "+err.Error())
		return
	}
	game.AddPlayer(id, req.PlayerName)

	WriteResponse(w, http.StatusCreated, map[string]any{
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

	fmt.Println(resp)
	WriteResponse(w, http.StatusOK, resp)
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
			WriteError(w, http.StatusNotFound, "Lobby not found")
			return
		}

		lobbyID, err := uuid.Parse(urlParts[1])
		if err != nil {
			WriteError(w, http.StatusBadRequest, "Invalid lobby ID")
			return
		}

		action := urlParts[2]
		switch action {
		case "join":
			joinLobby(w, r, lobbyID)
		default:
			WriteError(w, http.StatusNotImplemented, "Lobby action not implemented")
		}

	})

	mux.HandleFunc("GET /", defaultRoute)

	return mux
}
