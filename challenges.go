package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ChallengeType represents different types of daily challenges
type ChallengeType int

const (
	ChallengeSpeed    ChallengeType = iota // Reach target WPM
	ChallengeAccuracy                      // Maintain target accuracy
	ChallengeWords                         // Type target word count
	ChallengeTime                          // Type for target duration
	ChallengeNoErrors                      // Complete test with zero errors
	ChallengeStreak                        // Complete consecutive tests
)

// Challenge represents a daily challenge
type Challenge struct {
	ID          string        `json:"id"`
	Date        string        `json:"date"` // YYYY-MM-DD
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Type        ChallengeType `json:"type"`
	Target      int           `json:"target"` // Target value (WPM, accuracy %, words, etc.)
	Completed   bool          `json:"completed"`
	Progress    int           `json:"progress"` // Current progress toward target
	Reward      string        `json:"reward"`   // Description of reward
}

// DailyChallenges manages daily challenges
type DailyChallenges struct {
	Challenges []Challenge `json:"challenges"`
	path       string
}

// NewDailyChallenges creates or loads daily challenges
func NewDailyChallenges() *DailyChallenges {
	dc := &DailyChallenges{
		Challenges: []Challenge{},
	}

	// Get config directory
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = os.TempDir()
	}

	ktypeDir := filepath.Join(configDir, "ktype")
	if err := os.MkdirAll(ktypeDir, 0755); err != nil {
		dc.path = "challenges.json"
	} else {
		dc.path = filepath.Join(ktypeDir, "challenges.json")
	}

	dc.load()
	dc.ensureDailyChallenges()
	return dc
}

// load reads challenges from file
func (dc *DailyChallenges) load() {
	data, err := os.ReadFile(dc.path)
	if err != nil {
		return // File doesn't exist yet
	}

	if err := json.Unmarshal(data, dc); err != nil {
		dc.Challenges = []Challenge{}
	}
}

// save writes challenges to file
func (dc *DailyChallenges) save() error {
	data, err := json.MarshalIndent(dc, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(dc.path, data, 0644)
}

// ensureDailyChallenges creates challenges for today if they don't exist
func (dc *DailyChallenges) ensureDailyChallenges() {
	today := time.Now().Format("2006-01-02")

	// Check if we already have challenges for today
	hasToday := false
	for _, c := range dc.Challenges {
		if c.Date == today {
			hasToday = true
			break
		}
	}

	// Remove old challenges (keep only last 7 days)
	cutoff := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	var filtered []Challenge
	for _, c := range dc.Challenges {
		if c.Date >= cutoff {
			filtered = append(filtered, c)
		}
	}
	dc.Challenges = filtered

	// Generate today's challenges if needed
	if !hasToday {
		dc.generateDailyChallenges(today)
		dc.save()
	}
}

// generateDailyChallenges creates 3 random challenges for the day
func (dc *DailyChallenges) generateDailyChallenges(date string) {
	// Challenge 1: Speed challenge
	dc.Challenges = append(dc.Challenges, Challenge{
		ID:          fmt.Sprintf("%s-speed", date),
		Date:        date,
		Title:       "Speed Demon",
		Description: "Type at 60 WPM or higher in any mode",
		Type:        ChallengeSpeed,
		Target:      60,
		Completed:   false,
		Progress:    0,
		Reward:      "Speed Badge",
	})

	// Challenge 2: Accuracy challenge
	dc.Challenges = append(dc.Challenges, Challenge{
		ID:          fmt.Sprintf("%s-accuracy", date),
		Date:        date,
		Title:       "Perfectionist",
		Description: "Complete a test with 95% accuracy or higher",
		Type:        ChallengeAccuracy,
		Target:      95,
		Completed:   false,
		Progress:    0,
		Reward:      "Accuracy Badge",
	})

	// Challenge 3: Volume challenge
	dc.Challenges = append(dc.Challenges, Challenge{
		ID:          fmt.Sprintf("%s-volume", date),
		Date:        date,
		Title:       "Marathon Typist",
		Description: "Type 500 words total today",
		Type:        ChallengeWords,
		Target:      500,
		Completed:   false,
		Progress:    0,
		Reward:      "Endurance Badge",
	})
}

// GetTodaysChallenges returns today's challenges
func (dc *DailyChallenges) GetTodaysChallenges() []Challenge {
	today := time.Now().Format("2006-01-02")
	var todayChallenges []Challenge
	for _, c := range dc.Challenges {
		if c.Date == today {
			todayChallenges = append(todayChallenges, c)
		}
	}
	return todayChallenges
}

// UpdateProgress updates challenge progress based on a completed test
func (dc *DailyChallenges) UpdateProgress(wpm, accuracy, wordsTyped int) error {
	today := time.Now().Format("2006-01-02")

	for i := range dc.Challenges {
		if dc.Challenges[i].Date != today {
			continue
		}

		switch dc.Challenges[i].Type {
		case ChallengeSpeed:
			if wpm >= dc.Challenges[i].Target && !dc.Challenges[i].Completed {
				dc.Challenges[i].Progress = wpm
				dc.Challenges[i].Completed = true
			}

		case ChallengeAccuracy:
			if accuracy >= dc.Challenges[i].Target && !dc.Challenges[i].Completed {
				dc.Challenges[i].Progress = accuracy
				dc.Challenges[i].Completed = true
			}

		case ChallengeWords:
			dc.Challenges[i].Progress += wordsTyped
			if dc.Challenges[i].Progress >= dc.Challenges[i].Target {
				dc.Challenges[i].Completed = true
			}

		case ChallengeNoErrors:
			if accuracy == 100 && !dc.Challenges[i].Completed {
				dc.Challenges[i].Progress = 100
				dc.Challenges[i].Completed = true
			}
		}
	}

	return dc.save()
}

// GetCompletedCount returns number of completed challenges today
func (dc *DailyChallenges) GetCompletedCount() int {
	count := 0
	for _, c := range dc.GetTodaysChallenges() {
		if c.Completed {
			count++
		}
	}
	return count
}

// GetTotalCount returns total number of challenges today
func (dc *DailyChallenges) GetTotalCount() int {
	return len(dc.GetTodaysChallenges())
}

// HasCompletedAll returns true if all today's challenges are completed
func (dc *DailyChallenges) HasCompletedAll() bool {
	challenges := dc.GetTodaysChallenges()
	if len(challenges) == 0 {
		return false
	}

	for _, c := range challenges {
		if !c.Completed {
			return false
		}
	}
	return true
}
