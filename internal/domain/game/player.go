package game

import (
	"crypto/rand"
	"fmt"
	"sync"
)

type PlayerID [8]byte

type Player struct {
	ID         PlayerID
	Name       string
	Ready      bool
	WinAmount  int
	LoseAmount int
	WinRate    float32
	Hand       []Meld
}

var (
	Players   map[PlayerID]*Player
	PlayersMu sync.RWMutex
)

func NewPlayer(playerID PlayerID, playerName string) *Player {
	// Ensure only one routine updates player in-memory map at a time
	PlayersMu.Lock()
	defer PlayersMu.Unlock()

	player := &Player{
		ID:         playerID,
		Name:       playerName,
		WinAmount:  0,
		LoseAmount: 0,
		WinRate:    0,
		Hand:       nil,
	}
	Players[player.ID] = player

	return player
}

func GeneratePlayerId() (PlayerID, error) {
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var id PlayerID

	randomBytes := make([]byte, 8)
	if _, err := rand.Read(randomBytes); err != nil {
		return id, err
	}

	for i, b := range randomBytes {
		id[i] = charset[b%64]
	}

	return id, nil
}

func validateRawPlayerID(rawPlayerId string) error {
	var id PlayerID

	if len(rawPlayerId) != len(id) {
		return fmt.Errorf("invalid player_id")
	}

	return nil
}

func ConvertRawPlayerID(rawPlayerID string) (PlayerID, error) {
	var id PlayerID

	if err := validateRawPlayerID(rawPlayerID); err != nil {
		return id, err
	}

	copy(id[:], rawPlayerID)
	return id, nil
}
