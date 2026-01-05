package game

type Player struct {
	ID         int
	Name       string
	Ready      bool
	WinAmount  int
	LoseAmount int
	WinRate    float32
	Hand       []Meld
}

func NewPlayer(playerId int, playerName string) *Player {
	return &Player{
		ID:         playerId,
		Name:       playerName,
		WinAmount:  0,
		LoseAmount: 0,
		WinRate:    0,
		Hand:       nil,
	}
}
