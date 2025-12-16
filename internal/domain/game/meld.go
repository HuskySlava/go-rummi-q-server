package game

type Meld []Tile

const MinMeldLength = 3

func (tiles Meld) isAscendingSeries() bool {

	// Validate length
	if len(tiles) < MinMeldLength {
		return false
	}

	// Meld series should be a series of consecutive numbers
	for i := 0; i+1 < len(tiles); i++ {
		isJoker := tiles[i].IsJoker() || tiles[i+1].IsJoker()
		if !isJoker && tiles[i].Value+1 != tiles[i+1].Value {
			return false
		}
	}
	return true
}

func (tiles Meld) isUniqueColors() bool {

	// Map of existing colors
	seen := make(map[Color]bool)

	for _, v := range tiles {
		if !v.IsJoker() {
			if seen[v.Color] {
				return false
			}
			seen[v.Color] = true
		}
	}

	return true
}

func (tiles Meld) isGroup() bool {

	// Validate length
	if len(tiles) < MinMeldLength {
		return false
	}

	var expectedValue int
	var checkFromIndex int

	// Find first none Joker value and index in Meld
	foundNoneJoker := false
	for i, v := range tiles {
		if !v.IsJoker() {
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
	for _, t := range tiles[checkFromIndex:] {
		if t.Value != expectedValue && !t.IsJoker() {
			return false
		}
	}

	// Proceed with color check
	return tiles.isUniqueColors()
}

func (tiles Meld) IsValid() bool {
	return tiles.isAscendingSeries() || tiles.isGroup()
}
