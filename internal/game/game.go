package game

import (
	"fmt"
	"strings"
	"time"

	"ktype/internal/storage"
	"ktype/internal/words"
)

// Game holds all game state
type Game struct {
	Words        []string
	CurrentInput string
	WordIndex    int
	Correct      []bool
	TypedWords   []string

	StartTime time.Time
	Duration  time.Duration
	Elapsed   time.Duration

	Mode        Mode
	Difficulty  words.Difficulty
	Complexity  words.Complexity
	TargetWords int
	State       State
	TotalChars  int
	ErrorChars  int

	Heatmap *storage.Heatmap

	Errors        []TypingError
	CurrentErrors []int
}

// NewTimed creates a new timed game
func NewTimed(duration time.Duration, difficulty words.Difficulty, complexity words.Complexity, heatmap *storage.Heatmap) *Game {
	w := words.GetRandomWithComplexity(200, difficulty, complexity)
	return &Game{
		Words:         w,
		Correct:       make([]bool, 0),
		TypedWords:    make([]string, 0),
		Duration:      duration,
		Mode:          ModeTimed,
		Difficulty:    difficulty,
		Complexity:    complexity,
		State:         StatePlaying,
		Heatmap:       heatmap,
		Errors:        make([]TypingError, 0),
		CurrentErrors: make([]int, 0),
	}
}

// NewWords creates a new word-count game
func NewWords(wordCount int, difficulty words.Difficulty, complexity words.Complexity, heatmap *storage.Heatmap) *Game {
	w := words.GetRandomWithComplexity(wordCount+10, difficulty, complexity)
	return &Game{
		Words:         w,
		Correct:       make([]bool, 0),
		TypedWords:    make([]string, 0),
		TargetWords:   wordCount,
		Mode:          ModeWords,
		Difficulty:    difficulty,
		Complexity:    complexity,
		State:         StatePlaying,
		Heatmap:       heatmap,
		Errors:        make([]TypingError, 0),
		CurrentErrors: make([]int, 0),
	}
}

// NewZen creates a new zen mode game
func NewZen(difficulty words.Difficulty, complexity words.Complexity, heatmap *storage.Heatmap) *Game {
	w := words.GetRandomWithComplexity(1000, difficulty, complexity)
	return &Game{
		Words:         w,
		Correct:       make([]bool, 0),
		TypedWords:    make([]string, 0),
		Mode:          ModeZen,
		Difficulty:    difficulty,
		Complexity:    complexity,
		State:         StatePlaying,
		Heatmap:       heatmap,
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

	if g.Mode == ModeTimed {
		if g.Elapsed >= g.Duration {
			g.Elapsed = g.Duration
			g.State = StateFinished
		}
	}
}

// TimeRemaining returns the time remaining in seconds
func (g *Game) TimeRemaining() int {
	if g.Mode != ModeTimed {
		return -1
	}
	remaining := g.Duration - g.Elapsed
	if remaining < 0 {
		return 0
	}
	return int(remaining.Seconds())
}

// WordsRemaining returns words left to type
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

	g.Start()

	charStr := string(char)
	g.CurrentInput += charStr
	g.TotalChars++

	currentWord := g.Words[g.WordIndex]
	inputLen := len(g.CurrentInput)
	isCorrect := true

	if inputLen <= len(currentWord) {
		if g.CurrentInput[inputLen-1] != currentWord[inputLen-1] {
			g.ErrorChars++
			isCorrect = false
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
		g.ErrorChars++
		isCorrect = false
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

	g.Heatmap.RecordHit(charStr)
	if !isCorrect {
		g.Heatmap.RecordError(charStr)
	}
}

// HandleSpace processes space (move to next word)
func (g *Game) HandleSpace() {
	if g.State != StatePlaying {
		return
	}

	g.Start()

	currentWord := g.Words[g.WordIndex]
	isCorrect := g.CurrentInput == currentWord

	g.Correct = append(g.Correct, isCorrect)
	g.TypedWords = append(g.TypedWords, g.CurrentInput)
	g.TotalChars++

	g.CurrentInput = ""
	g.WordIndex++

	if g.WordIndex >= len(g.Words) {
		if g.Mode == ModeZen {
			newWords := words.GetRandomWithComplexity(100, g.Difficulty, g.Complexity)
			g.Words = append(g.Words, newWords...)
		} else {
			g.State = StateFinished
			return
		}
	}

	if g.Mode == ModeWords && len(g.TypedWords) >= g.TargetWords {
		g.State = StateFinished
	}
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
func (g *Game) WPM() int {
	if g.Elapsed == 0 {
		return 0
	}

	correctChars := 0
	for i, typed := range g.TypedWords {
		if i < len(g.Words) && typed == g.Words[i] {
			correctChars += len(typed) + 1
		}
	}

	minutes := g.Elapsed.Minutes()
	if minutes == 0 {
		return 0
	}

	return int(float64(correctChars) / 5.0 / minutes)
}

// RawWPM calculates raw WPM
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

// CurrentWordState returns the state of the current word
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
