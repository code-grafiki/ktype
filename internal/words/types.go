package words

// Difficulty represents word difficulty level
type Difficulty int

const (
	DifficultyEasy Difficulty = iota
	DifficultyMedium
	DifficultyHard
)

// String returns a string representation of difficulty
func (d Difficulty) String() string {
	switch d {
	case DifficultyEasy:
		return "easy"
	case DifficultyMedium:
		return "medium"
	case DifficultyHard:
		return "hard"
	default:
		return "medium"
	}
}

// Complexity represents additional complexity options
type Complexity int

const (
	ComplexityNormal Complexity = iota
	ComplexityPunctuation
	ComplexityNumbers
	ComplexityFull
)

// String returns a string representation
func (c Complexity) String() string {
	switch c {
	case ComplexityPunctuation:
		return "punctuation"
	case ComplexityNumbers:
		return "numbers"
	case ComplexityFull:
		return "full"
	default:
		return "normal"
	}
}
