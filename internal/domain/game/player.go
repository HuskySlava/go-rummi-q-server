package game

import (
	"crypto/rand"
	"fmt"
	"strings"
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
	Players   = make(map[PlayerID]*Player)
	PlayersMu sync.RWMutex
)

const playerIDCharset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

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
		Hand:       []Meld{},
	}
	Players[player.ID] = player

	return player
}

func GeneratePlayerID() (PlayerID, error) {
	var id PlayerID

	randomBytes := make([]byte, 8)
	if _, err := rand.Read(randomBytes); err != nil {
		return id, err
	}

	for i, b := range randomBytes {
		id[i] = playerIDCharset[b%byte(len(playerIDCharset))]
	}

	return id, nil
}

func validateRawPlayerID(rawPlayerId string) error {
	var id PlayerID

	if len(rawPlayerId) != len(id) {
		return fmt.Errorf("invalid player_id")
	}

	for _, c := range rawPlayerId {
		if !strings.ContainsRune(playerIDCharset, c) {
			return fmt.Errorf("invalid player_id character: %q", c)
		}
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
