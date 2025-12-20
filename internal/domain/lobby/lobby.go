package lobby

import (
	"github.com/google/uuid"
	"go-rummi-q-server/internal/domain/game"
	"time"
)

type Lobby struct {
	ID         uuid.UUID
	StartTime  time.Time
	LastActive time.Time
	IsGame     bool
	Players    *[]game.Player
	Game       *game.Game
}

func NewLobby() *Lobby {
	return &Lobby{
		ID:         uuid.New(),
		StartTime:  time.Now(),
		LastActive: time.Now(),
		IsGame:     false,
	}
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
