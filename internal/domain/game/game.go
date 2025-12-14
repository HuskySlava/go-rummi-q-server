package game

type Player struct {
	ID         int
	Name       string
	WinAmount  int
	LoseAmount int
	WinRate    float32
	Hand       []Meld
}

type Game struct {
	GameID     int
	Players    []Player
	Board      []Meld
	PlayerTurn Player
	PlayerWon  Player
}

func (g *Game) NewGame() {
	// TODO
}

func (g *Game) NextTurn() {
	// TODO
}
