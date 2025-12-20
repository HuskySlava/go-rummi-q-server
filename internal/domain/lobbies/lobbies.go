package lobbies

import (
	"github.com/google/uuid"
	"go-rummi-q-server/internal/domain/game"
	"sync"
	"time"
)

var (
	lobbies   = make(map[uuid.UUID]*Lobby)
	lobbiesMu sync.RWMutex
)

type Lobby struct {
	mu         sync.RWMutex
	ID         uuid.UUID
	StartTime  time.Time
	LastActive time.Time
	IsGame     bool
	Players    []game.Player
	Game       *game.Game
}

func NewLobby() *Lobby {
	lobby := &Lobby{
		ID:         uuid.New(),
		StartTime:  time.Now(),
		LastActive: time.Now(),
		IsGame:     false,
		Players:    make([]game.Player, 0),
	}

	// Ensure only one routine updates lobby in-memory map at a time
	lobbiesMu.Lock()
	lobbies[lobby.ID] = lobby
	lobbiesMu.Unlock()

	return lobby
}

func LobbyExists(id uuid.UUID) bool {
	// Block writing to lobbies while you read from it
	lobbiesMu.RLock()
	defer lobbiesMu.RUnlock()

	_, ok := lobbies[id]
	return ok
}

func TerminateLobby(id uuid.UUID) {
	// Block writing to lobbies while it is being deleted
	lobbiesMu.Lock()
	defer lobbiesMu.Unlock()
	delete(lobbies, id)
}

func NewPlayer() *game.Player {
	return &game.Player{
		ID:         0,
		Name:       "",
		WinAmount:  0,
		LoseAmount: 0,
		WinRate:    0,
		Hand:       nil,
	}
}

func PlayerJoin() {

}
