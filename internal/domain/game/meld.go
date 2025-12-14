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

func (tiles Meld) isGroup() bool {

	// Validate length
	if len(tiles) < MinMeldLength {
		return false
	}

	var expectedValue uint8
	var checkFromIndex int

	// Find first none Joker value and index in Meld
	foundNoneJoker := false
	for i, v := range tiles {
		if v.Value != JokerValue {
			expectedValue = v.Value
			checkFromIndex = i
			foundNoneJoker = true
			break
		}
	}

	// All jokers (Impossible on classic rules, future-proof for custom rules)
	if !foundNoneJoker {
		return true
	}

	// All tiles in a Meld group must be of the same value
	for i := checkFromIndex; i < len(tiles); i++ { // Skip irrelevant indexes (jokers
		if tiles[i].Value != expectedValue && tiles[i].Value != JokerValue {
			return false
		}
	}

	// Proceed with color check
	return true
}

func (tiles Meld) IsValid() bool {
	return tiles.isAscendingSeries() || tiles.isGroup()
}
