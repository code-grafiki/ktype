package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewLeaderboard(t *testing.T) {
	lb := NewLeaderboard()

	if lb == nil {
		t.Fatal("NewLeaderboard returned nil")
	}

	if lb.Scores == nil {
		t.Error("Scores slice should be initialized")
	}

	if lb.path == "" {
		t.Error("Path should be set")
	}
}

func TestLeaderboardAddScore(t *testing.T) {
	// Create temp directory for test
	tempDir := t.TempDir()

	lb := &Leaderboard{
		Scores: []Score{},
		path:   filepath.Join(tempDir, "test_scores.json"),
	}

	lb.AddScore(65, 95, "time:30")

	if len(lb.Scores) != 1 {
		t.Errorf("Expected 1 score, got %d", len(lb.Scores))
	}

	score := lb.Scores[0]
	if score.WPM != 65 {
		t.Errorf("Expected WPM 65, got %d", score.WPM)
	}

	if score.Accuracy != 95 {
		t.Errorf("Expected accuracy 95, got %d", score.Accuracy)
	}

	if score.Mode != "time:30" {
		t.Errorf("Expected mode 'time:30', got %q", score.Mode)
	}

	if score.Date.IsZero() {
		t.Error("Expected date to be set")
	}
}

func TestLeaderboardGetPB(t *testing.T) {
	lb := &Leaderboard{
		Scores: []Score{
			{WPM: 50, Accuracy: 90, Mode: "time:30", Date: time.Now()},
			{WPM: 65, Accuracy: 92, Mode: "time:30", Date: time.Now()},
			{WPM: 45, Accuracy: 95, Mode: "time:30", Date: time.Now()},
			{WPM: 70, Accuracy: 88, Mode: "time:60", Date: time.Now()},
		},
	}

	// Get PB for time:30
	pb := lb.GetPB("time:30")
	if pb == nil {
		t.Fatal("Expected PB to be found")
	}
	if pb.WPM != 65 {
		t.Errorf("Expected PB WPM 65, got %d", pb.WPM)
	}

	// Get PB for time:60
	pb = lb.GetPB("time:60")
	if pb == nil || pb.WPM != 70 {
		t.Errorf("Expected PB WPM 70 for time:60, got %v", pb)
	}

	// Get overall PB
	pb = lb.GetOverallPB()
	if pb == nil || pb.WPM != 70 {
		t.Errorf("Expected overall PB WPM 70, got %v", pb)
	}

	// Get PB for mode with no scores
	pb = lb.GetPB("words:100")
	if pb != nil {
		t.Error("Expected nil PB for mode with no scores")
	}
}

func TestLeaderboardIsPB(t *testing.T) {
	lb := &Leaderboard{
		Scores: []Score{
			{WPM: 60, Accuracy: 90, Mode: "time:30", Date: time.Now()},
		},
	}

	if !lb.IsPB(70, "time:30") {
		t.Error("70 WPM should be new PB")
	}

	if lb.IsPB(50, "time:30") {
		t.Error("50 WPM should not be PB when 60 exists")
	}

	// First score for a mode should always be PB
	if !lb.IsPB(30, "words:25") {
		t.Error("First score for a mode should be PB")
	}
}

func TestLeaderboardGetTopScores(t *testing.T) {
	lb := &Leaderboard{
		Scores: []Score{
			{WPM: 50, Accuracy: 90, Mode: "time:30", Date: time.Now()},
			{WPM: 65, Accuracy: 92, Mode: "time:30", Date: time.Now()},
			{WPM: 45, Accuracy: 95, Mode: "time:30", Date: time.Now()},
			{WPM: 70, Accuracy: 88, Mode: "time:30", Date: time.Now()},
			{WPM: 55, Accuracy: 91, Mode: "time:60", Date: time.Now()},
		},
	}

	top := lb.GetTopScores("time:30", 3)
	if len(top) != 3 {
		t.Errorf("Expected 3 top scores, got %d", len(top))
	}

	// Should be sorted by WPM descending
	if top[0].WPM != 70 || top[1].WPM != 65 || top[2].WPM != 50 {
		t.Error("Top scores not sorted correctly")
	}

	// Request more than available
	top = lb.GetTopScores("time:60", 10)
	if len(top) != 1 {
		t.Errorf("Expected 1 score, got %d", len(top))
	}
}

func TestLeaderboardSaveAndLoad(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "scores.json")

	// Create and save
	lb1 := &Leaderboard{
		Scores: []Score{
			{WPM: 65, Accuracy: 95, Mode: "time:30", Date: time.Now()},
		},
		path: path,
	}

	if err := lb1.save(); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("Score file should exist after save")
	}

	// Load into new leaderboard
	lb2 := &Leaderboard{
		Scores: []Score{},
		path:   path,
	}
	lb2.load()

	if len(lb2.Scores) != 1 {
		t.Errorf("Expected 1 score after load, got %d", len(lb2.Scores))
	}

	if lb2.Scores[0].WPM != 65 {
		t.Errorf("Expected WPM 65 after load, got %d", lb2.Scores[0].WPM)
	}
}

func TestLeaderboardLoadCorruptFile(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "corrupt.json")

	// Write invalid JSON
	if err := os.WriteFile(path, []byte("not valid json"), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	lb := &Leaderboard{
		Scores: []Score{{WPM: 100, Accuracy: 99, Mode: "test", Date: time.Now()}},
		path:   path,
	}
	lb.load()

	// Should reset to empty on corrupt file
	if len(lb.Scores) != 0 {
		t.Errorf("Expected empty scores after loading corrupt file, got %d", len(lb.Scores))
	}
}

func TestLeaderboardLoadNonExistentFile(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "nonexistent.json")

	lb := &Leaderboard{
		Scores: []Score{},
		path:   path,
	}
	lb.load()

	// Should not error, just leave scores empty
	if len(lb.Scores) != 0 {
		t.Errorf("Expected empty scores, got %d", len(lb.Scores))
	}
}

func TestLeaderboardScoreLimit(t *testing.T) {
	tempDir := t.TempDir()
	lb := &Leaderboard{
		Scores: []Score{},
		path:   filepath.Join(tempDir, "scores.json"),
	}

	// Add 105 scores
	for i := 0; i < 105; i++ {
		lb.AddScore(i, 90, "time:30")
	}

	if len(lb.Scores) != 100 {
		t.Errorf("Expected 100 scores (limit), got %d", len(lb.Scores))
	}

	// Verify it kept the most recent scores
	if lb.Scores[0].WPM != 5 {
		t.Error("Should have trimmed oldest scores, keeping most recent")
	}
}

func TestScoreJSONSerialization(t *testing.T) {
	score := Score{
		WPM:      65,
		Accuracy: 95,
		Mode:     "time:30",
		Date:     time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
	}

	data, err := json.Marshal(score)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var decoded Score
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if decoded.WPM != score.WPM {
		t.Errorf("WPM mismatch: got %d, want %d", decoded.WPM, score.WPM)
	}

	if decoded.Accuracy != score.Accuracy {
		t.Errorf("Accuracy mismatch: got %d, want %d", decoded.Accuracy, score.Accuracy)
	}

	if decoded.Mode != score.Mode {
		t.Errorf("Mode mismatch: got %q, want %q", decoded.Mode, score.Mode)
	}
}
