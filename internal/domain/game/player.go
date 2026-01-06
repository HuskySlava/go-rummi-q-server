package game

import (
	"crypto/rand"
	"encoding/json"
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
	players   = make(map[PlayerID]*Player)
	playersMu sync.RWMutex
)

const playerIDCharset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func AddPlayer(playerID PlayerID, playerName string) *Player {
	// Ensure only one routine updates player in-memory map at a time
	playersMu.Lock()
	defer playersMu.Unlock()

	player := &Player{
		ID:         playerID,
		Name:       playerName,
		WinAmount:  0,
		LoseAmount: 0,
		WinRate:    0,
		Hand:       []Meld{},
	}
	players[player.ID] = player

	return player
}

func GetAllPlayersJSON() ([]byte, error) {
	playersMu.RLock()
	defer playersMu.RUnlock()

	var result []map[string]string

	for k, v := range players {
		result = append(result, map[string]string{
			"id":    string(k[:]),
			"value": v.Name,
		})
	}

	data := map[string]interface{}{
		"players": result,
	}

	return json.MarshalIndent(data, "", "  ")
}

func GetPlayer(playerID PlayerID) (*Player, error) {
	player, ok := players[playerID]
	if !ok {
		return nil, fmt.Errorf("player %v not found", playerID)
	}
	return player, nil
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

	if len(rawPlayerId) != len(PlayerID{}) {
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
