package main

import (
	"fmt"
	"strings"
	"time"
)

// GameState represents the current state of the game
type GameState int

const (
	StateMenu             GameState = iota
	StateModeSelect                 // Choose timed vs words
	StateTimeSelect                 // Choose duration
	StateWordsSelect                // Choose word count
	StateCustomInput                // Custom input
	StateDifficultySelect           // Choose difficulty
	StateComplexitySelect           // Choose word complexity
	StateStats                      // Statistics dashboard
	StateHeatmap                    // Typing heatmap
	StateCustomWordList             // Custom word lists
	StatePlaying
	StateFinished
)

// GameMode represents the type of game
type GameMode int

const (
	ModeTimed GameMode = iota
	ModeWords
	ModeZen
)

// TypingError represents a single typing error with details
type TypingError struct {
	ExpectedChar rune      // The character that should have been typed
	TypedChar    rune      // The character that was actually typed
	Position     int       // Position in the word where error occurred
	WordIndex    int       // Which word the error occurred in
	Timestamp    time.Time // When the error occurred
	ErrorType    ErrorType // Type of error
}

// ErrorType categorizes different types of typing errors
type ErrorType int

const (
	ErrorWrongChar     ErrorType = iota // Typed wrong character
	ErrorExtraChar                      // Typed extra character beyond word length
	ErrorMissingChar                    // Didn't type a character (skipped)
	ErrorTransposition                  // Swapped adjacent characters
)

// Game holds all game state
type Game struct {
	Words        []string // Words to type
	CurrentInput string   // Current word being typed
	WordIndex    int      // Current word position
	Correct      []bool   // Track if each word was typed correctly
	TypedWords   []string // What the user actually typed

	StartTime time.Time
	Duration  time.Duration
	Elapsed   time.Duration

	Mode        GameMode
	Difficulty  Difficulty
	Complexity  WordComplexity
	TargetWords int // For words mode
	State       GameState
	TotalChars  int // Total characters typed
	ErrorChars  int // Characters typed incorrectly

	heatmap *Heatmap // Model-level heatmap reference

	// Enhanced error tracking
	Errors        []TypingError // Detailed error log
	CurrentErrors []int         // Positions of errors in current word
}

// NewTimedGame creates a new timed game
func NewTimedGame(duration time.Duration, difficulty Difficulty, complexity WordComplexity, heatmap *Heatmap) *Game {
	words := getRandomWordsWithComplexity(200, difficulty, complexity) // Get enough words for any test
	return &Game{
		Words:         words,
		Correct:       make([]bool, 0),
		TypedWords:    make([]string, 0),
		Duration:      duration,
		Mode:          ModeTimed,
		Difficulty:    difficulty,
		Complexity:    complexity,
		State:         StatePlaying,
		heatmap:       heatmap,
		Errors:        make([]TypingError, 0),
		CurrentErrors: make([]int, 0),
	}
}

// NewWordsGame creates a new word-count game
func NewWordsGame(wordCount int, difficulty Difficulty, complexity WordComplexity, heatmap *Heatmap) *Game {
	words := getRandomWordsWithComplexity(wordCount+10, difficulty, complexity) // A few extra just in case
	return &Game{
		Words:         words,
		Correct:       make([]bool, 0),
		TypedWords:    make([]string, 0),
		TargetWords:   wordCount,
		Mode:          ModeWords,
		Difficulty:    difficulty,
		Complexity:    complexity,
		State:         StatePlaying,
		heatmap:       heatmap,
		Errors:        make([]TypingError, 0),
		CurrentErrors: make([]int, 0),
	}
}

// NewZenGame creates a new zen mode game (unlimited typing)
func NewZenGame(difficulty Difficulty, complexity WordComplexity, heatmap *Heatmap) *Game {
	words := getRandomWordsWithComplexity(1000, difficulty, complexity) // Large pool for zen mode
	return &Game{
		Words:         words,
		Correct:       make([]bool, 0),
		TypedWords:    make([]string, 0),
		Mode:          ModeZen,
		Difficulty:    difficulty,
		Complexity:    complexity,
		State:         StatePlaying,
		heatmap:       heatmap,
		Errors:        make([]TypingError, 0),
		CurrentErrors: make([]int, 0),
	}
}

// ModeString returns a string representation for leaderboard
func (g *Game) ModeString() string {
	if g.Mode == ModeTimed {
		return fmt.Sprintf("time:%d", int(g.Duration.Seconds()))
	}
	if g.Mode == ModeWords {
		return fmt.Sprintf("words:%d", g.TargetWords)
	}
	return "zen"
}

// Start begins the game timer
func (g *Game) Start() {
	if g.StartTime.IsZero() {
		g.StartTime = time.Now()
	}
}

// Update updates the elapsed time and checks if game is finished
func (g *Game) Update() {
	if g.State != StatePlaying || g.StartTime.IsZero() {
		return
	}

	g.Elapsed = time.Since(g.StartTime)

	// Check finish condition based on mode
	if g.Mode == ModeTimed {
		if g.Elapsed >= g.Duration {
			g.Elapsed = g.Duration
			g.State = StateFinished
		}
	}
	// Words mode finishes when target reached (handled in HandleSpace)
	// Zen mode never finishes automatically
}

// TimeRemaining returns the time remaining in seconds (for timed mode)
func (g *Game) TimeRemaining() int {
	if g.Mode != ModeTimed {
		return -1 // Indicates N/A
	}
	remaining := g.Duration - g.Elapsed
	if remaining < 0 {
		return 0
	}
	return int(remaining.Seconds())
}

// WordsRemaining returns words left to type (for words mode)
func (g *Game) WordsRemaining() int {
	if g.Mode != ModeWords {
		return -1
	}
	remaining := g.TargetWords - len(g.TypedWords)
	if remaining < 0 {
		return 0
	}
	return remaining
}

// Progress returns progress string for display
func (g *Game) Progress() string {
	if g.Mode == ModeTimed {
		return fmt.Sprintf("%ds", g.TimeRemaining())
	}
	if g.Mode == ModeWords {
		return fmt.Sprintf("%d/%d", len(g.TypedWords), g.TargetWords)
	}
	return fmt.Sprintf("%d words", len(g.TypedWords))
}

// HandleChar processes a typed character
func (g *Game) HandleChar(char rune) {
	if g.State != StatePlaying {
		return
	}

	// Start timer on first keystroke
	g.Start()

	charStr := string(char)
	g.CurrentInput += charStr
	g.TotalChars++

	// Check if current character is correct
	currentWord := g.Words[g.WordIndex]
	inputLen := len(g.CurrentInput)
	isCorrect := true

	if inputLen <= len(currentWord) {
		if g.CurrentInput[inputLen-1] != currentWord[inputLen-1] {
			g.ErrorChars++
			isCorrect = false
			// Record detailed error
			err := TypingError{
				ExpectedChar: rune(currentWord[inputLen-1]),
				TypedChar:    char,
				Position:     inputLen - 1,
				WordIndex:    g.WordIndex,
				Timestamp:    time.Now(),
				ErrorType:    ErrorWrongChar,
			}
			g.Errors = append(g.Errors, err)
			g.CurrentErrors = append(g.CurrentErrors, inputLen-1)
		}
	} else {
		g.ErrorChars++ // Typed more chars than the word has
		isCorrect = false
		// Record detailed error for extra character
		err := TypingError{
			ExpectedChar: 0,
			TypedChar:    char,
			Position:     inputLen - 1,
			WordIndex:    g.WordIndex,
			Timestamp:    time.Now(),
			ErrorType:    ErrorExtraChar,
		}
		g.Errors = append(g.Errors, err)
		g.CurrentErrors = append(g.CurrentErrors, inputLen-1)
	}

	// Track in heatmap
	g.heatmap.RecordHit(charStr)
	if !isCorrect {
		g.heatmap.RecordError(charStr)
	}
}

// HandleSpace processes space (move to next word)
func (g *Game) HandleSpace() {
	if g.State != StatePlaying {
		return
	}

	// Start timer on first keystroke
	g.Start()

	currentWord := g.Words[g.WordIndex]
	isCorrect := g.CurrentInput == currentWord

	g.Correct = append(g.Correct, isCorrect)
	g.TypedWords = append(g.TypedWords, g.CurrentInput)
	g.TotalChars++ // Count space

	g.CurrentInput = ""
	g.WordIndex++

	// Check finish conditions
	if g.WordIndex >= len(g.Words) {
		// For zen mode, generate more words if we run out
		if g.Mode == ModeZen {
			// Add more words to the pool
			newWords := getRandomWordsWithComplexity(100, g.Difficulty, g.Complexity)
			g.Words = append(g.Words, newWords...)
		} else {
			g.State = StateFinished
			return
		}
	}

	// Words mode: check if target reached
	if g.Mode == ModeWords && len(g.TypedWords) >= g.TargetWords {
		g.State = StateFinished
	}
	// Zen mode never finishes
}

// HandleBackspace removes the last character
func (g *Game) HandleBackspace() {
	if g.State != StatePlaying {
		return
	}

	if len(g.CurrentInput) > 0 {
		g.CurrentInput = g.CurrentInput[:len(g.CurrentInput)-1]
	}
}

// WPM calculates words per minute
// Standard: 5 characters = 1 word
func (g *Game) WPM() int {
	if g.Elapsed == 0 {
		return 0
	}

	// Count total correct characters
	correctChars := 0
	for i, typed := range g.TypedWords {
		if i < len(g.Words) && typed == g.Words[i] {
			correctChars += len(typed) + 1 // +1 for space
		}
	}

	minutes := g.Elapsed.Minutes()
	if minutes == 0 {
		return 0
	}

	// WPM = (correct chars / 5) / minutes
	return int(float64(correctChars) / 5.0 / minutes)
}

// RawWPM calculates raw WPM (all typed characters, including errors)
func (g *Game) RawWPM() int {
	if g.Elapsed == 0 {
		return 0
	}

	minutes := g.Elapsed.Minutes()
	if minutes == 0 {
		return 0
	}

	return int(float64(g.TotalChars) / 5.0 / minutes)
}

// Accuracy returns the accuracy percentage
func (g *Game) Accuracy() int {
	if g.TotalChars == 0 {
		return 100
	}

	correctChars := g.TotalChars - g.ErrorChars
	if correctChars < 0 {
		correctChars = 0
	}

	return int(float64(correctChars) / float64(g.TotalChars) * 100)
}

// CurrentWordState returns the state of the current word being typed
// Returns: (correct part, error part, remaining part)
func (g *Game) CurrentWordState() (string, string, string) {
	if g.WordIndex >= len(g.Words) {
		return "", "", ""
	}

	word := g.Words[g.WordIndex]
	input := g.CurrentInput

	var correct, errors, remaining strings.Builder

	for i := 0; i < len(word); i++ {
		if i < len(input) {
			if input[i] == word[i] {
				correct.WriteByte(word[i])
			} else {
				errors.WriteByte(word[i])
			}
		} else {
			remaining.WriteByte(word[i])
		}
	}

	// Extra characters typed beyond word length
	if len(input) > len(word) {
		errors.WriteString(input[len(word):])
	}

	return correct.String(), errors.String(), remaining.String()
}

// CorrectWordsCount returns the number of correctly typed words
func (g *Game) CorrectWordsCount() int {
	count := 0
	for _, c := range g.Correct {
		if c {
			count++
		}
	}
	return count
}
