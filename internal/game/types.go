package game

import "time"

// State represents the current state of the application
type State int

const (
	StateMenu State = iota
	StateModeSelect
	StateTimeSelect
	StateWordsSelect
	StateCustomInput
	StateDifficultySelect
	StateComplexitySelect
	StateStats
	StateHeatmap
	StateCustomWordList
	StateSettings
	StateCursorSelect
	StateColorSelect
	StateChallenges
	StatePlaying
	StateFinished
)

// Mode represents the type of game
type Mode int

const (
	ModeTimed Mode = iota
	ModeWords
	ModeZen
)

// ErrorType categorizes different types of typing errors
type ErrorType int

const (
	ErrorWrongChar ErrorType = iota
	ErrorExtraChar
	ErrorMissingChar
	ErrorTransposition
)

// TypingError represents a single typing error with details
type TypingError struct {
	ExpectedChar rune
	TypedChar    rune
	Position     int
	WordIndex    int
	Timestamp    time.Time
	ErrorType    ErrorType
}
