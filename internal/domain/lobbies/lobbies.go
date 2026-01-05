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
	ReadyToStart
	GameInProgress
	GameEnded
)

type Lobby struct {
	mu           sync.RWMutex
	ID           uuid.UUID
	StartTime    time.Time
	LastActive   time.Time
	Status       LobbyStatus
	Players      []game.Player
	Game         *game.Game
	NextPlayerID int
}

// ## Constructors ##

func NewLobby() *Lobby {
	// Ensure only one routine updates lobby in-memory map at a time

	lobbiesMu.Lock()
	defer lobbiesMu.Unlock()

	lobby := &Lobby{
		ID:           uuid.New(),
		StartTime:    time.Now(),
		LastActive:   time.Now(),
		Players:      make([]game.Player, 0),
		NextPlayerID: 1,
		Status:       AwaitingPlayers,
	}

	lobbies[lobby.ID] = lobby

	return lobby
}

// ## Methods ##

func (l *Lobby) join(playerName string, rawPlayerID string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	playerId, err := game.ConvertRawPlayerID(rawPlayerID)
	if err != nil {
		return err
	}

	player := game.NewPlayer(playerId, playerName)
	l.Players = append(l.Players, *player)
	l.NextPlayerID++

	// More than 1 player joined, game is ready to start
	if len(l.Players) > 1 {
		var err error
		l.Game, err = game.NewGame(l.Players)
		if err != nil {
			return err
		}
		l.Status = ReadyToStart
	}
	l.LastActive = time.Now()

	return nil
}

// ## helpers ##

func LobbyExists(id uuid.UUID) bool {
	// Block writing to lobbies while you read from it
	lobbiesMu.RLock()
	defer lobbiesMu.RUnlock()

	_, ok := lobbies[id]
	return ok
}

func JoinLobby(lobbyId uuid.UUID, playerID string, playerName string) error {
	lobbiesMu.RLock() // Protect lobbies map
	lobby, ok := lobbies[lobbyId]
	lobbiesMu.RUnlock()

	if !ok {
		return fmt.Errorf("failed to join lobby")
	}

	err := lobby.join(playerName, playerID)
	if err != nil {
		return err
	}

	return nil
}

func TerminateLobby(id uuid.UUID) {
	// Block writing to lobbies while it is being deleted
	lobbiesMu.Lock()
	defer lobbiesMu.Unlock()
	delete(lobbies, id)

	// TODO: Signal lobby is terminated
}
