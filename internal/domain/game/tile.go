package game

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
	Value uint8
	Color Color
}

func NewTile(color Color, value uint8) Tile {
	return Tile{
		ID:    string(color) + "-" + string(value),
		Value: value,
		Color: color,
	}
}

func (t Tile) IsJoker() bool {
	return t.Value == JokerValue
}
