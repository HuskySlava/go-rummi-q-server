package game

import (
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

type Game struct {
	GameID          uuid.UUID
	TilePool        []Tile
	Board           []Meld
	Players         []Player
	PlayerTurnIndex int
	PlayerWon       *Player
}

func NewGame(players []Player) (*Game, error) {

	var g = &Game{}

	if len(players) == 0 {
		err := fmt.Errorf("cannot start game without players")
		return nil, err
	}

	g.GameID = uuid.New()

	// Generate and shuffle tiles
	g.TilePool = generateTilePool()
	g.shuffleTiles()

	g.Board = make([]Meld, 0)
	g.Players = players

	g.dealTiles()

	// Set initial player turn
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	g.PlayerTurnIndex = r.Intn(len(g.Players))

	return g, nil
}

func (g *Game) NextTurn() {
	if len(g.Players) == 0 {
		return
	}
	if g.PlayerTurnIndex+1 >= len(g.Players) {
		g.PlayerTurnIndex = 0
	} else {
		g.PlayerTurnIndex++
	}
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

func (g *Game) shuffleTiles() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(g.TilePool), func(i, j int) {
		g.TilePool[i], g.TilePool[j] = g.TilePool[j], g.TilePool[i]
	})
}

func (g *Game) dealTiles() {
	const initialTilesPerPlayer = 14
	playersAmount := len(g.Players)

	for i := 0; i < playersAmount; i++ {
		start := i * initialTilesPerPlayer
		end := start + initialTilesPerPlayer
		playerTiles := []Meld{g.TilePool[start:end]}
		g.Players[i].Hand = append([]Meld{}, playerTiles...)
	}

	// Remove from the pool
	g.TilePool = g.TilePool[playersAmount*initialTilesPerPlayer:]
}
