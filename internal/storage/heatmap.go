package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// KeyStats tracks statistics for individual keys
type KeyStats struct {
	Key        string    `json:"key"`
	TotalHits  int       `json:"total_hits"`
	ErrorCount int       `json:"error_count"`
	LastUsed   time.Time `json:"last_used"`
}

// ErrorRate returns the error percentage for this key
func (k *KeyStats) ErrorRate() float64 {
	if k.TotalHits == 0 {
		return 0.0
	}
	return float64(k.ErrorCount) / float64(k.TotalHits) * 100.0
}

// Heatmap stores keystroke statistics for all keys
type Heatmap struct {
	Keys map[string]*KeyStats `json:"keys"`
	path string
}

// NewHeatmap creates or loads a heatmap
func NewHeatmap() *Heatmap {
	h := &Heatmap{
		Keys: make(map[string]*KeyStats),
	}

	// Get config directory
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = os.TempDir()
	}

	ktypeDir := filepath.Join(configDir, "ktype")
	if err := os.MkdirAll(ktypeDir, 0755); err != nil {
		h.path = "heatmap.json"
	} else {
		h.path = filepath.Join(ktypeDir, "heatmap.json")
	}

	h.load()
	return h
}

// load reads heatmap from file
func (h *Heatmap) load() {
	data, err := os.ReadFile(h.path)
	if err != nil {
		return // File doesn't exist yet
	}

	if err := json.Unmarshal(data, h); err != nil {
		// Corrupt file - start fresh
		h.Keys = make(map[string]*KeyStats)
	}
}

// save writes heatmap to file
func (h *Heatmap) save() error {
	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(h.path, data, 0644)
}

// RecordHit records a successful keystroke
func (h *Heatmap) RecordHit(key string) {
	if key == "" {
		return
	}

	// Normalize key
	key = strings.ToLower(key)
	if len(key) > 1 {
		// Handle special keys
		key = normalizeKey(key)
	}

	if _, exists := h.Keys[key]; !exists {
		h.Keys[key] = &KeyStats{Key: key}
	}

	h.Keys[key].TotalHits++
	h.Keys[key].LastUsed = time.Now()
	h.save()
}

// RecordError records an error for a key
func (h *Heatmap) RecordError(key string) {
	if key == "" {
		return
	}

	// Normalize key
	key = strings.ToLower(key)
	if len(key) > 1 {
		key = normalizeKey(key)
	}

	if _, exists := h.Keys[key]; !exists {
		h.Keys[key] = &KeyStats{Key: key}
	}

	h.Keys[key].ErrorCount++
	h.Keys[key].LastUsed = time.Now()
	h.save()
}

// normalizeKey converts special key names to standard format
func normalizeKey(key string) string {
	// Map special keys to simpler representations
	switch key {
	case "space":
		return " "
	case "enter", "return":
		return "↵"
	case "backspace", "delete":
		return "⌫"
	case "tab":
		return "⇥"
	case "esc":
		return "esc"
	default:
		// For other multi-character keys, just return first char or the key itself
		if len(key) == 1 {
			return key
		}
		return key[:1]
	}
}

// GetTopErrors returns keys with the most errors, sorted by error rate
func (h *Heatmap) GetTopErrors(limit int) []*KeyStats {
	var stats []*KeyStats
	for _, stat := range h.Keys {
		if stat.ErrorCount > 0 {
			stats = append(stats, stat)
		}
	}

	// Sort by error rate (descending)
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].ErrorRate() > stats[j].ErrorRate()
	})

	if len(stats) > limit {
		stats = stats[:limit]
	}

	return stats
}

// GetMostUsed returns the most frequently typed keys
func (h *Heatmap) GetMostUsed(limit int) []*KeyStats {
	var stats []*KeyStats
	for _, stat := range h.Keys {
		stats = append(stats, stat)
	}

	// Sort by total hits (descending)
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].TotalHits > stats[j].TotalHits
	})

	if len(stats) > limit {
		stats = stats[:limit]
	}

	return stats
}

// GetHeatmapData returns heatmap data organized by keyboard rows
func (h *Heatmap) GetHeatmapData() KeyboardHeatmap {
	return KeyboardHeatmap{
		TopRow:    h.getRowStats("qwertyuiop"),
		HomeRow:   h.getRowStats("asdfghjkl;"),
		BottomRow: h.getRowStats("zxcvbnm,./"),
		Numbers:   h.getRowStats("1234567890"),
	}
}

// getRowStats gets stats for a specific row of keys
func (h *Heatmap) getRowStats(keys string) []*KeyStats {
	var stats []*KeyStats
	for _, key := range keys {
		keyStr := string(key)
		if stat, exists := h.Keys[keyStr]; exists {
			stats = append(stats, stat)
		} else {
			// Return empty stat for keys not yet typed
			stats = append(stats, &KeyStats{Key: keyStr})
		}
	}
	return stats
}

// KeyboardHeatmap organizes heatmap data by keyboard rows
type KeyboardHeatmap struct {
	TopRow    []*KeyStats
	HomeRow   []*KeyStats
	BottomRow []*KeyStats
	Numbers   []*KeyStats
}

// GetErrorHeatLevel returns a heat level (0-4) based on error rate
func GetErrorHeatLevel(errorRate float64) int {
	switch {
	case errorRate == 0:
		return 0 // No errors
	case errorRate < 5:
		return 1 // Low errors
	case errorRate < 15:
		return 2 // Medium errors
	case errorRate < 30:
		return 3 // High errors
	default:
		return 4 // Very high errors
	}
}

// GetHeatColor returns a color for a heat level
func GetHeatColor(level int) string {
	colors := []string{
		"#646669", // Level 0: Gray (no data/no errors)
		"#98c379", // Level 1: Green (low errors)
		"#e2b714", // Level 2: Yellow (medium errors)
		"#d19a66", // Level 3: Orange (high errors)
		"#ca4754", // Level 4: Red (very high errors)
	}
	if level < 0 || level >= len(colors) {
		return colors[0]
	}
	return colors[level]
}

// GetTotalKeystrokes returns total keystrokes recorded
func (h *Heatmap) GetTotalKeystrokes() int {
	total := 0
	for _, stat := range h.Keys {
		total += stat.TotalHits
	}
	return total
}

// GetTotalErrors returns total errors recorded
func (h *Heatmap) GetTotalErrors() int {
	total := 0
	for _, stat := range h.Keys {
		total += stat.ErrorCount
	}
	return total
}

// GetOverallAccuracy returns overall accuracy percentage
func (h *Heatmap) GetOverallAccuracy() float64 {
	totalHits := h.GetTotalKeystrokes()
	totalErrors := h.GetTotalErrors()

	if totalHits == 0 {
		return 100.0
	}

	return float64(totalHits-totalErrors) / float64(totalHits) * 100.0
}

// Clear resets all heatmap data
func (h *Heatmap) Clear() {
	h.Keys = make(map[string]*KeyStats)
	h.save()
}
