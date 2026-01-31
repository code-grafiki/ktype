package game

import (
	"testing"
	"time"

	"ktype/internal/storage"
	"ktype/internal/words"
)

func TestNewTimed(t *testing.T) {
	game := NewTimed(30*time.Second, words.DifficultyMedium, words.ComplexityNormal, storage.NewHeatmap())

	if game == nil {
		t.Fatal("NewTimed returned nil")
	}

	if game.Mode != ModeTimed {
		t.Errorf("Expected ModeTimed, got %v", game.Mode)
	}

	if game.Duration != 30*time.Second {
		t.Errorf("Expected 30s duration, got %v", game.Duration)
	}

	if game.State != StatePlaying {
		t.Errorf("Expected StatePlaying, got %v", game.State)
	}

	if len(game.Words) == 0 {
		t.Error("Expected words to be populated")
	}
}

func TestNewWords(t *testing.T) {
	game := NewWords(25, words.DifficultyMedium, words.ComplexityNormal, storage.NewHeatmap())

	if game == nil {
		t.Fatal("NewWords returned nil")
	}

	if game.Mode != ModeWords {
		t.Errorf("Expected ModeWords, got %v", game.Mode)
	}

	if game.TargetWords != 25 {
		t.Errorf("Expected 25 target words, got %d", game.TargetWords)
	}

	if len(game.Words) < 25 {
		t.Errorf("Expected at least 25 words, got %d", len(game.Words))
	}
}

func TestNewZen(t *testing.T) {
	game := NewZen(words.DifficultyMedium, words.ComplexityNormal, storage.NewHeatmap())

	if game == nil {
		t.Fatal("NewZen returned nil")
	}

	if game.Mode != ModeZen {
		t.Errorf("Expected ModeZen, got %v", game.Mode)
	}

	if game.Duration != 0 {
		t.Errorf("Expected no duration limit for zen mode, got %v", game.Duration)
	}
}

func TestModeString(t *testing.T) {
	tests := []struct {
		mode     Mode
		duration time.Duration
		target   int
		expected string
	}{
		{ModeTimed, 30 * time.Second, 0, "time:30"},
		{ModeTimed, 60 * time.Second, 0, "time:60"},
		{ModeWords, 0, 25, "words:25"},
		{ModeWords, 0, 50, "words:50"},
		{ModeZen, 0, 0, "zen"},
	}

	for _, tt := range tests {
		var g *Game
		if tt.mode == ModeTimed {
			g = NewTimed(tt.duration, words.DifficultyMedium, words.ComplexityNormal, storage.NewHeatmap())
		} else if tt.mode == ModeWords {
			g = NewWords(tt.target, words.DifficultyMedium, words.ComplexityNormal, storage.NewHeatmap())
		} else {
			g = NewZen(words.DifficultyMedium, words.ComplexityNormal, storage.NewHeatmap())
		}

		result := g.ModeString()
		if result != tt.expected {
			t.Errorf("ModeString() = %q, want %q", result, tt.expected)
		}
	}
}

func TestWPM(t *testing.T) {
	g := NewTimed(30*time.Second, words.DifficultyMedium, words.ComplexityNormal, storage.NewHeatmap())

	// Test WPM at start (should be 0)
	if wpm := g.WPM(); wpm != 0 {
		t.Errorf("Expected WPM 0 at start, got %d", wpm)
	}

	// Simulate typing
	g.StartTime = time.Now().Add(-1 * time.Minute) // 1 minute ago
	g.TypedWords = []string{"hello", "world"}
	g.Correct = []bool{true, true}
	g.Words = []string{"hello", "world", "test"}
	g.Elapsed = 1 * time.Minute // Set elapsed directly

	// WPM should be calculated based on correct words
	// (5+1) + (5+1) = 12 chars / 5 = 2.4 words / 1 minute = 2.4 WPM ≈ 2
	wpm := g.WPM()
	if wpm <= 0 {
		t.Errorf("Expected positive WPM, got %d", wpm)
	}
}

func TestAccuracy(t *testing.T) {
	tests := []struct {
		name       string
		totalChars int
		errorChars int
		expected   int
	}{
		{"perfect", 100, 0, 100},
		{"half errors", 100, 50, 50},
		{"all errors", 100, 100, 0},
		{"no input", 0, 0, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewTimed(30*time.Second, words.DifficultyMedium, words.ComplexityNormal, storage.NewHeatmap())
			g.TotalChars = tt.totalChars
			g.ErrorChars = tt.errorChars

			acc := g.Accuracy()
			if acc != tt.expected {
				t.Errorf("Accuracy() = %d, want %d", acc, tt.expected)
			}
		})
	}
}

func TestHandleChar(t *testing.T) {
	g := NewTimed(30*time.Second, words.DifficultyMedium, words.ComplexityNormal, storage.NewHeatmap())
	g.Words = []string{"hello", "world"}

	// Type "he" correctly
	g.HandleChar('h')
	g.HandleChar('e')

	if g.CurrentInput != "he" {
		t.Errorf("Expected 'he', got %q", g.CurrentInput)
	}

	if g.TotalChars != 2 {
		t.Errorf("Expected 2 total chars, got %d", g.TotalChars)
	}

	// Type 'x' (error)
	g.HandleChar('x')
	if g.ErrorChars != 1 {
		t.Errorf("Expected 1 error char, got %d", g.ErrorChars)
	}

	// Timer should have started
	if g.StartTime.IsZero() {
		t.Error("Expected timer to start on first keystroke")
	}
}

func TestHandleSpace(t *testing.T) {
	g := NewWords(2, words.DifficultyMedium, words.ComplexityNormal, storage.NewHeatmap())
	g.Words = []string{"hello", "world", "test"}
	g.CurrentInput = "hello"

	g.HandleSpace()

	if g.WordIndex != 1 {
		t.Errorf("Expected WordIndex 1, got %d", g.WordIndex)
	}

	if len(g.TypedWords) != 1 {
		t.Errorf("Expected 1 typed word, got %d", len(g.TypedWords))
	}

	if g.CurrentInput != "" {
		t.Errorf("Expected CurrentInput to be cleared, got %q", g.CurrentInput)
	}
}

func TestHandleBackspace(t *testing.T) {
	g := NewTimed(30*time.Second, words.DifficultyMedium, words.ComplexityNormal, storage.NewHeatmap())
	g.CurrentInput = "hello"

	g.HandleBackspace()
	if g.CurrentInput != "hell" {
		t.Errorf("Expected 'hell', got %q", g.CurrentInput)
	}

	// Backspace on empty input should not crash
	g.CurrentInput = ""
	g.HandleBackspace()
	if g.CurrentInput != "" {
		t.Errorf("Expected empty string, got %q", g.CurrentInput)
	}
}

func TestTimeRemaining(t *testing.T) {
	g := NewTimed(30*time.Second, words.DifficultyMedium, words.ComplexityNormal, storage.NewHeatmap())
	g.StartTime = time.Now()
	g.Elapsed = 10 * time.Second

	remaining := g.TimeRemaining()
	if remaining != 20 {
		t.Errorf("Expected 20s remaining, got %d", remaining)
	}

	// Words mode should return -1
	wordsGame := NewWords(25, words.DifficultyMedium, words.ComplexityNormal, storage.NewHeatmap())
	if wordsGame.TimeRemaining() != -1 {
		t.Error("Expected -1 for words mode")
	}
}

func TestWordsRemaining(t *testing.T) {
	g := NewWords(10, words.DifficultyMedium, words.ComplexityNormal, storage.NewHeatmap())
	g.TypedWords = make([]string, 3)

	remaining := g.WordsRemaining()
	if remaining != 7 {
		t.Errorf("Expected 7 words remaining, got %d", remaining)
	}

	// Timed mode should return -1
	timedGame := NewTimed(30*time.Second, words.DifficultyMedium, words.ComplexityNormal, storage.NewHeatmap())
	if timedGame.WordsRemaining() != -1 {
		t.Error("Expected -1 for timed mode")
	}
}

func TestCorrectWordsCount(t *testing.T) {
	g := NewTimed(30*time.Second, words.DifficultyMedium, words.ComplexityNormal, storage.NewHeatmap())
	g.Correct = []bool{true, false, true, true}

	count := g.CorrectWordsCount()
	if count != 3 {
		t.Errorf("Expected 3 correct words, got %d", count)
	}
}

func TestCurrentWordState(t *testing.T) {
	g := NewTimed(30*time.Second, words.DifficultyMedium, words.ComplexityNormal, storage.NewHeatmap())
	g.Words = []string{"hello"}
	g.CurrentInput = "he"

	correct, errors, remaining := g.CurrentWordState()

	if correct != "he" {
		t.Errorf("Expected correct 'he', got %q", correct)
	}

	if errors != "" {
		t.Errorf("Expected no errors, got %q", errors)
	}

	if remaining != "llo" {
		t.Errorf("Expected remaining 'llo', got %q", remaining)
	}
}

func TestUpdate(t *testing.T) {
	g := NewTimed(1*time.Second, words.DifficultyMedium, words.ComplexityNormal, storage.NewHeatmap())
	g.StartTime = time.Now().Add(-2 * time.Second) // Started 2s ago

	g.Update()

	if g.State != StateFinished {
		t.Error("Expected game to finish when time expires")
	}

	if g.Elapsed > g.Duration {
		t.Error("Elapsed time should not exceed duration")
	}
}
