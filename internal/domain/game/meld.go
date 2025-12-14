package game

type Meld []Tile

const MinMeldLength = 3
const JokerValue = 99

func (tiles Meld) isAscendingSeries() bool {

	// Validate length
	if len(tiles) < MinMeldLength {
		return false
	}

	// Meld series should be a series of consecutive numbers
	for i := 0; i+1 < len(tiles); i++ {
		isJoker := tiles[i].Value == JokerValue || tiles[i+1].Value == JokerValue
		if !isJoker && tiles[i].Value+1 != tiles[i+1].Value {
			return false
		}
	}
	return true
}

func (tiles Meld) IsValid() bool {
	return tiles.isAscendingSeries()
}
