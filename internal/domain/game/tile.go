package game

type Color string

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
		ID:    "R1-1",
		Value: 1,
		Color: ColorPurple,
	}
}
