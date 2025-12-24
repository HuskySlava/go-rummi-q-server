package game

import (
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

type Player struct {
	ID         int
	Name       string
	WinAmount  int
	LoseAmount int
	WinRate    float32
	Hand       []Meld
}

type Game struct {
	GameID     uuid.UUID
	TilePool   []Tile
	Board      []Meld
	Players    []Player
	PlayerTurn *Player
	PlayerWon  *Player
}

func NewGame(players []Player) (*Game, error) {

	var g = &Game{}

	if len(players) == 0 {
		err := fmt.Errorf("cannot start game without players")
		return nil, err
	}

	g.GameID = uuid.New()
	g.TilePool = generateTilePool()
	g.Board = make([]Meld, 0)
	g.Players = players

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	firstPlayerIndex := r.Intn(len(g.Players))
	g.PlayerTurn = &g.Players[firstPlayerIndex]

	return g, nil
}

func (g *Game) NextTurn() {
	// TODO
}

// ## helpers ##

func generateTiles(amount int, color Color, copies int) []Tile {
	tiles := make([]Tile, 0, amount*copies)
	for i := 1; i <= amount; i++ {
		for c := 0; c < copies; c++ {
			if color == ColorJoker {
				tiles = append(tiles, NewTile(color, JokerValue))
			} else {
				tiles = append(tiles, NewTile(color, i))
			}
		}
	}
	return tiles
}

func generateTilePool() []Tile {

	const StandardTileCount = 13
	const TileCopies = 2

	const JokerTileCount = 1
	const JokerCopies = 2

	type colorInfo struct {
		amount int
		copies int
		color  Color
	}

	colors := []colorInfo{
		{amount: StandardTileCount, color: ColorRed, copies: TileCopies},
		{amount: StandardTileCount, color: ColorBlue, copies: TileCopies},
		{amount: StandardTileCount, color: ColorPurple, copies: TileCopies},
		{amount: StandardTileCount, color: ColorBlack, copies: TileCopies},
		{amount: JokerTileCount, color: ColorJoker, copies: JokerCopies},
	}

	poolSize := 0
	for _, ci := range colors {
		poolSize += ci.amount * ci.copies
	}

	tiles := make([]Tile, 0, poolSize)

	for _, ci := range colors {
		tiles = append(tiles, generateTiles(ci.amount, ci.color, ci.copies)...)
	}

	return tiles
}
