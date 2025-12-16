package game

import "strconv"

type Color string

const JokerValue = 99

const (
	ColorRed    Color = "RED"
	ColorBlue   Color = "BLUE"
	ColorBlack  Color = "BLACK"
	ColorPurple Color = "PURPLE"
	ColorJoker  Color = "JOKER"
)

type Tile struct {
	ID    string
	Value int
	Color Color
}

func NewTile(color Color, value int) Tile {
	return Tile{
		ID:    string(color) + "-" + strconv.Itoa(value),
		Value: value,
		Color: color,
	}
}

func (t Tile) IsJoker() bool {
	return t.Value == JokerValue
}
