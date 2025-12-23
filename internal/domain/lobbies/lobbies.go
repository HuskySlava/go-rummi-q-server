package lobbies

import (
	"fmt"
	"github.com/google/uuid"
	"go-rummi-q-server/internal/domain/game"
	"sync"
	"time"
)

var (
	lobbies   = make(map[uuid.UUID]*Lobby)
	lobbiesMu sync.RWMutex
)

type LobbyStatus int

const (
	AwaitingPlayers LobbyStatus = iota + 1
	GameInProgress
	GameEnded
)

type Lobby struct {
	mu           sync.RWMutex
	ID           uuid.UUID
	StartTime    time.Time
	LastActive   time.Time
	Status       LobbyStatus
	Players      []*game.Player
	Game         *game.Game
	NextPlayerID int
}

func NewLobby() *Lobby {
	lobbiesMu.Lock()
	defer lobbiesMu.Unlock()

	lobby := &Lobby{
		ID:           uuid.New(),
		StartTime:    time.Now(),
		LastActive:   time.Now(),
		Players:      make([]*game.Player, 0),
		NextPlayerID: 1,
		Status:       AwaitingPlayers,
	}

	// Ensure only one routine updates lobby in-memory map at a time
	lobbies[lobby.ID] = lobby

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

	// TODO: Signal lobby is terminated
}

func NewPlayer(playerId int, playerName string) *game.Player {
	return &game.Player{
		ID:         playerId,
		Name:       playerName,
		WinAmount:  0,
		LoseAmount: 0,
		WinRate:    0,
		Hand:       nil,
	}
}

func (l *Lobby) join(playerName string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	player := NewPlayer(l.NextPlayerID, playerName)
	l.Players = append(l.Players, player)
	l.NextPlayerID++

	if len(l.Players) > 1 {
		l.Status = GameInProgress
	}
}

func JoinLobby(lobbyId uuid.UUID, playerName string) error {
	lobbiesMu.RLock() // Protect lobbies map
	lobby, ok := lobbies[lobbyId]
	lobbiesMu.RUnlock()
	if !ok {
		return fmt.Errorf("failed to join lobby")
	}
	lobby.join(playerName)
	return nil
}
