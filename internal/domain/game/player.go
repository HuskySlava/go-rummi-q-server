package game

import "crypto/rand"

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

var Players map[PlayerID]*Player

func NewPlayer(playerID PlayerID, playerName string) *Player {
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
