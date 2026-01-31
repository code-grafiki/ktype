package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Score represents a single test result
type Score struct {
	WPM      int       `json:"wpm"`
	Accuracy int       `json:"accuracy"`
	Mode     string    `json:"mode"` // "time:30", "time:60", "words:25", etc.
	Date     time.Time `json:"date"`
}

// Leaderboard manages local scores
type Leaderboard struct {
	Scores []Score `json:"scores"`
	path   string
}

// NewLeaderboard creates or loads a leaderboard
func NewLeaderboard() *Leaderboard {
	lb := &Leaderboard{
		Scores: []Score{},
	}

	// Get config directory
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = os.TempDir()
	}

	ktypeDir := filepath.Join(configDir, "ktype")
	if err := os.MkdirAll(ktypeDir, 0755); err != nil {
		// Fall back to current directory
		lb.path = "scores.json"
	} else {
		lb.path = filepath.Join(ktypeDir, "scores.json")
	}

	lb.load()
	return lb
}

// load reads scores from file
func (lb *Leaderboard) load() {
	data, err := os.ReadFile(lb.path)
	if err != nil {
		return // File doesn't exist yet
	}

	if err := json.Unmarshal(data, lb); err != nil {
		// Corrupt file - start fresh
		lb.Scores = []Score{}
	}
}

// save writes scores to file
func (lb *Leaderboard) save() error {
	data, err := json.MarshalIndent(lb, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(lb.path, data, 0644)
}

// AddScore adds a new score and saves
func (lb *Leaderboard) AddScore(wpm, accuracy int, mode string) {
	score := Score{
		WPM:      wpm,
		Accuracy: accuracy,
		Mode:     mode,
		Date:     time.Now(),
	}

	lb.Scores = append(lb.Scores, score)

	// Keep only last 100 scores
	if len(lb.Scores) > 100 {
		lb.Scores = lb.Scores[len(lb.Scores)-100:]
	}

	lb.save()
}

// GetPB returns the personal best WPM for a mode (or overall if mode is empty)
func (lb *Leaderboard) GetPB(mode string) *Score {
	var best *Score

	for i := range lb.Scores {
		score := &lb.Scores[i]

		// Filter by mode if specified
		if mode != "" && score.Mode != mode {
			continue
		}

		if best == nil || score.WPM > best.WPM {
			best = score
		}
	}

	return best
}

// GetOverallPB returns the best score across all modes
func (lb *Leaderboard) GetOverallPB() *Score {
	return lb.GetPB("")
}

// GetTopScores returns top N scores for a mode
func (lb *Leaderboard) GetTopScores(mode string, n int) []Score {
	var filtered []Score

	for _, score := range lb.Scores {
		if mode == "" || score.Mode == mode {
			filtered = append(filtered, score)
		}
	}

	// Sort by WPM descending
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].WPM > filtered[j].WPM
	})

	if len(filtered) > n {
		filtered = filtered[:n]
	}

	return filtered
}

// IsPB checks if a WPM would be a new personal best for a mode
func (lb *Leaderboard) IsPB(wpm int, mode string) bool {
	pb := lb.GetPB(mode)
	return pb == nil || wpm > pb.WPM
}
