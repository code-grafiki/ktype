package main

import (
	"testing"
	"time"
)

func TestNewTimedGame(t *testing.T) {
	game := NewTimedGame(30*time.Second, DifficultyMedium, ComplexityNormal, NewHeatmap())

	if game == nil {
		t.Fatal("NewTimedGame returned nil")
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

func TestNewWordsGame(t *testing.T) {
	game := NewWordsGame(25, DifficultyMedium, ComplexityNormal, NewHeatmap())

	if game == nil {
		t.Fatal("NewWordsGame returned nil")
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

func TestNewZenGame(t *testing.T) {
	game := NewZenGame(DifficultyMedium, ComplexityNormal, NewHeatmap())

	if game == nil {
		t.Fatal("NewZenGame returned nil")
	}

	if game.Mode != ModeZen {
		t.Errorf("Expected ModeZen, got %v", game.Mode)
	}

	if game.Duration != 0 {
		t.Errorf("Expected no duration limit for zen mode, got %v", game.Duration)
	}
}

func TestGameModeString(t *testing.T) {
	tests := []struct {
		mode     GameMode
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
		var game *Game
		if tt.mode == ModeTimed {
			game = NewTimedGame(tt.duration, DifficultyMedium, ComplexityNormal, NewHeatmap())
		} else if tt.mode == ModeWords {
			game = NewWordsGame(tt.target, DifficultyMedium, ComplexityNormal, NewHeatmap())
		} else {
			game = NewZenGame(DifficultyMedium, ComplexityNormal, NewHeatmap())
		}

		result := game.ModeString()
		if result != tt.expected {
			t.Errorf("ModeString() = %q, want %q", result, tt.expected)
		}
	}
}

func TestGameWPM(t *testing.T) {
	game := NewTimedGame(30*time.Second, DifficultyMedium, ComplexityNormal, NewHeatmap())

	// Test WPM at start (should be 0)
	if wpm := game.WPM(); wpm != 0 {
		t.Errorf("Expected WPM 0 at start, got %d", wpm)
	}

	// Simulate typing
	game.StartTime = time.Now().Add(-1 * time.Minute) // 1 minute ago
	game.TypedWords = []string{"hello", "world"}
	game.Correct = []bool{true, true}
	game.Words = []string{"hello", "world", "test"}
	game.Elapsed = 1 * time.Minute // Set elapsed directly

	// WPM should be calculated based on correct words
	// (5+1) + (5+1) = 12 chars / 5 = 2.4 words / 1 minute = 2.4 WPM ≈ 2
	wpm := game.WPM()
	if wpm <= 0 {
		t.Errorf("Expected positive WPM, got %d", wpm)
	}
}

func TestGameAccuracy(t *testing.T) {
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
			game := NewTimedGame(30*time.Second, DifficultyMedium, ComplexityNormal, NewHeatmap())
			game.TotalChars = tt.totalChars
			game.ErrorChars = tt.errorChars

			acc := game.Accuracy()
			if acc != tt.expected {
				t.Errorf("Accuracy() = %d, want %d", acc, tt.expected)
			}
		})
	}
}

func TestGameHandleChar(t *testing.T) {
	game := NewTimedGame(30*time.Second, DifficultyMedium, ComplexityNormal, NewHeatmap())
	game.Words = []string{"hello", "world"}

	// Type "he" correctly
	game.HandleChar('h')
	game.HandleChar('e')

	if game.CurrentInput != "he" {
		t.Errorf("Expected 'he', got %q", game.CurrentInput)
	}

	if game.TotalChars != 2 {
		t.Errorf("Expected 2 total chars, got %d", game.TotalChars)
	}

	// Type 'x' (error)
	game.HandleChar('x')
	if game.ErrorChars != 1 {
		t.Errorf("Expected 1 error char, got %d", game.ErrorChars)
	}

	// Timer should have started
	if game.StartTime.IsZero() {
		t.Error("Expected timer to start on first keystroke")
	}
}

func TestGameHandleSpace(t *testing.T) {
	game := NewWordsGame(2, DifficultyMedium, ComplexityNormal, NewHeatmap())
	game.Words = []string{"hello", "world", "test"}
	game.CurrentInput = "hello"

	game.HandleSpace()

	if game.WordIndex != 1 {
		t.Errorf("Expected WordIndex 1, got %d", game.WordIndex)
	}

	if len(game.TypedWords) != 1 {
		t.Errorf("Expected 1 typed word, got %d", len(game.TypedWords))
	}

	if game.CurrentInput != "" {
		t.Errorf("Expected CurrentInput to be cleared, got %q", game.CurrentInput)
	}
}

func TestGameHandleBackspace(t *testing.T) {
	game := NewTimedGame(30*time.Second, DifficultyMedium, ComplexityNormal, NewHeatmap())
	game.CurrentInput = "hello"

	game.HandleBackspace()
	if game.CurrentInput != "hell" {
		t.Errorf("Expected 'hell', got %q", game.CurrentInput)
	}

	// Backspace on empty input should not crash
	game.CurrentInput = ""
	game.HandleBackspace()
	if game.CurrentInput != "" {
		t.Errorf("Expected empty string, got %q", game.CurrentInput)
	}
}

func TestGameTimeRemaining(t *testing.T) {
	game := NewTimedGame(30*time.Second, DifficultyMedium, ComplexityNormal, NewHeatmap())
	game.StartTime = time.Now()
	game.Elapsed = 10 * time.Second

	remaining := game.TimeRemaining()
	if remaining != 20 {
		t.Errorf("Expected 20s remaining, got %d", remaining)
	}

	// Words mode should return -1
	wordsGame := NewWordsGame(25, DifficultyMedium, ComplexityNormal, NewHeatmap())
	if wordsGame.TimeRemaining() != -1 {
		t.Error("Expected -1 for words mode")
	}
}

func TestGameWordsRemaining(t *testing.T) {
	game := NewWordsGame(10, DifficultyMedium, ComplexityNormal, NewHeatmap())
	game.TypedWords = make([]string, 3)

	remaining := game.WordsRemaining()
	if remaining != 7 {
		t.Errorf("Expected 7 words remaining, got %d", remaining)
	}

	// Timed mode should return -1
	timedGame := NewTimedGame(30*time.Second, DifficultyMedium, ComplexityNormal, NewHeatmap())
	if timedGame.WordsRemaining() != -1 {
		t.Error("Expected -1 for timed mode")
	}
}

func TestGameCorrectWordsCount(t *testing.T) {
	game := NewTimedGame(30*time.Second, DifficultyMedium, ComplexityNormal, NewHeatmap())
	game.Correct = []bool{true, false, true, true}

	count := game.CorrectWordsCount()
	if count != 3 {
		t.Errorf("Expected 3 correct words, got %d", count)
	}
}

func TestGameCurrentWordState(t *testing.T) {
	game := NewTimedGame(30*time.Second, DifficultyMedium, ComplexityNormal, NewHeatmap())
	game.Words = []string{"hello"}
	game.CurrentInput = "he"

	correct, errors, remaining := game.CurrentWordState()

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

func TestGameUpdate(t *testing.T) {
	game := NewTimedGame(1*time.Second, DifficultyMedium, ComplexityNormal, NewHeatmap())
	game.StartTime = time.Now().Add(-2 * time.Second) // Started 2s ago

	game.Update()

	if game.State != StateFinished {
		t.Error("Expected game to finish when time expires")
	}

	if game.Elapsed > game.Duration {
		t.Error("Elapsed time should not exceed duration")
	}
}
