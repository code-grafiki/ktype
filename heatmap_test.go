package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewHeatmap(t *testing.T) {
	hm := NewHeatmap()
	if hm == nil {
		t.Fatal("NewHeatmap returned nil")
	}
	if hm.Keys == nil {
		t.Error("Heatmap.Keys should be initialized")
	}
}

func TestHeatmapRecordHit(t *testing.T) {
	hm := &Heatmap{
		Keys: make(map[string]*KeyStats),
		path: filepath.Join(os.TempDir(), "test_heatmap.json"),
	}

	hm.RecordHit("a")
	hm.RecordHit("a")
	hm.RecordHit("b")

	if hm.Keys["a"].TotalHits != 2 {
		t.Errorf("Expected 2 hits for 'a', got %d", hm.Keys["a"].TotalHits)
	}

	if hm.Keys["b"].TotalHits != 1 {
		t.Errorf("Expected 1 hit for 'b', got %d", hm.Keys["b"].TotalHits)
	}
}

func TestHeatmapRecordError(t *testing.T) {
	hm := &Heatmap{
		Keys: make(map[string]*KeyStats),
		path: filepath.Join(os.TempDir(), "test_heatmap.json"),
	}

	hm.RecordHit("a")
	hm.RecordHit("a")
	hm.RecordError("a")

	if hm.Keys["a"].ErrorCount != 1 {
		t.Errorf("Expected 1 error for 'a', got %d", hm.Keys["a"].ErrorCount)
	}

	if hm.Keys["a"].TotalHits != 2 {
		t.Errorf("Expected 2 hits for 'a', got %d", hm.Keys["a"].TotalHits)
	}
}

func TestKeyStatsErrorRate(t *testing.T) {
	ks := &KeyStats{
		Key:        "a",
		TotalHits:  10,
		ErrorCount: 2,
	}

	rate := ks.ErrorRate()
	expected := 20.0
	if rate != expected {
		t.Errorf("Expected error rate %.1f%%, got %.1f%%", expected, rate)
	}

	// Test with no hits
	ks2 := &KeyStats{Key: "b"}
	if ks2.ErrorRate() != 0.0 {
		t.Error("Error rate should be 0 with no hits")
	}
}

func TestNormalizeKey(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"a", "a"},
		{"A", "A"}, // Single char is returned as-is (strings.ToLower is called before normalizeKey)
		{"space", " "},
		{"enter", "↵"},
		{"return", "↵"},
		{"backspace", "⌫"},
		{"delete", "⌫"},
		{"tab", "⇥"},
		{"esc", "esc"},
		{"abc", "a"},
	}

	for _, tt := range tests {
		result := normalizeKey(tt.input)
		if result != tt.expected {
			t.Errorf("normalizeKey(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestGetTopErrors(t *testing.T) {
	hm := &Heatmap{
		Keys: map[string]*KeyStats{
			"a": {Key: "a", TotalHits: 10, ErrorCount: 5, LastUsed: time.Now()},
			"b": {Key: "b", TotalHits: 10, ErrorCount: 3, LastUsed: time.Now()},
			"c": {Key: "c", TotalHits: 10, ErrorCount: 7, LastUsed: time.Now()},
			"d": {Key: "d", TotalHits: 10, ErrorCount: 1, LastUsed: time.Now()},
		},
		path: filepath.Join(os.TempDir(), "test_heatmap.json"),
	}

	topErrors := hm.GetTopErrors(3)

	if len(topErrors) != 3 {
		t.Errorf("Expected 3 top errors, got %d", len(topErrors))
	}

	// Should be sorted by error rate descending
	if topErrors[0].Key != "c" {
		t.Errorf("Expected 'c' to have highest error rate, got %s", topErrors[0].Key)
	}

	if topErrors[1].Key != "a" {
		t.Errorf("Expected 'a' to have second highest error rate, got %s", topErrors[1].Key)
	}
}

func TestGetMostUsed(t *testing.T) {
	hm := &Heatmap{
		Keys: map[string]*KeyStats{
			"a": {Key: "a", TotalHits: 50, ErrorCount: 0, LastUsed: time.Now()},
			"b": {Key: "b", TotalHits: 30, ErrorCount: 0, LastUsed: time.Now()},
			"c": {Key: "c", TotalHits: 70, ErrorCount: 0, LastUsed: time.Now()},
		},
		path: filepath.Join(os.TempDir(), "test_heatmap.json"),
	}

	mostUsed := hm.GetMostUsed(2)

	if len(mostUsed) != 2 {
		t.Errorf("Expected 2 most used, got %d", len(mostUsed))
	}

	if mostUsed[0].Key != "c" {
		t.Errorf("Expected 'c' to be most used, got %s", mostUsed[0].Key)
	}
}

func TestGetHeatmapData(t *testing.T) {
	hm := &Heatmap{
		Keys: map[string]*KeyStats{
			"q": {Key: "q", TotalHits: 5, ErrorCount: 1, LastUsed: time.Now()},
			"a": {Key: "a", TotalHits: 10, ErrorCount: 2, LastUsed: time.Now()},
			"z": {Key: "z", TotalHits: 3, ErrorCount: 0, LastUsed: time.Now()},
			"1": {Key: "1", TotalHits: 2, ErrorCount: 1, LastUsed: time.Now()},
		},
		path: filepath.Join(os.TempDir(), "test_heatmap.json"),
	}

	data := hm.GetHeatmapData()

	// Check that we have the right number of keys in each row
	if len(data.TopRow) != 10 {
		t.Errorf("Expected 10 keys in top row, got %d", len(data.TopRow))
	}

	if len(data.HomeRow) != 10 {
		t.Errorf("Expected 10 keys in home row, got %d", len(data.HomeRow))
	}

	if len(data.BottomRow) != 10 {
		t.Errorf("Expected 10 keys in bottom row, got %d", len(data.BottomRow))
	}

	if len(data.Numbers) != 10 {
		t.Errorf("Expected 10 keys in numbers row, got %d", len(data.Numbers))
	}

	// Check that existing keys have data
	for _, stat := range data.TopRow {
		if stat.Key == "q" && stat.TotalHits != 5 {
			t.Errorf("Expected 'q' to have 5 hits, got %d", stat.TotalHits)
		}
	}

	// Check that missing keys have empty stats
	for _, stat := range data.TopRow {
		if stat.Key == "w" && stat.TotalHits != 0 {
			t.Errorf("Expected 'w' to have 0 hits (not typed), got %d", stat.TotalHits)
		}
	}
}

func TestGetErrorHeatLevel(t *testing.T) {
	tests := []struct {
		rate     float64
		expected int
	}{
		{0, 0},
		{3, 1},
		{10, 2},
		{20, 3},
		{40, 4},
		{100, 4},
	}

	for _, tt := range tests {
		result := GetErrorHeatLevel(tt.rate)
		if result != tt.expected {
			t.Errorf("GetErrorHeatLevel(%.1f) = %d, expected %d", tt.rate, result, tt.expected)
		}
	}
}

func TestGetHeatColor(t *testing.T) {
	// Test valid levels
	colors := []string{
		"#646669", // Level 0
		"#98c379", // Level 1
		"#e2b714", // Level 2
		"#d19a66", // Level 3
		"#ca4754", // Level 4
	}

	for i, expected := range colors {
		result := GetHeatColor(i)
		if result != expected {
			t.Errorf("GetHeatColor(%d) = %s, expected %s", i, result, expected)
		}
	}

	// Test invalid levels
	if GetHeatColor(-1) != colors[0] {
		t.Error("GetHeatColor(-1) should return default color")
	}

	if GetHeatColor(10) != colors[0] {
		t.Error("GetHeatColor(10) should return default color")
	}
}

func TestHeatmapTotals(t *testing.T) {
	hm := &Heatmap{
		Keys: map[string]*KeyStats{
			"a": {Key: "a", TotalHits: 100, ErrorCount: 10},
			"b": {Key: "b", TotalHits: 50, ErrorCount: 5},
			"c": {Key: "c", TotalHits: 25, ErrorCount: 0},
		},
		path: filepath.Join(os.TempDir(), "test_heatmap.json"),
	}

	if hm.GetTotalKeystrokes() != 175 {
		t.Errorf("Expected 175 total keystrokes, got %d", hm.GetTotalKeystrokes())
	}

	if hm.GetTotalErrors() != 15 {
		t.Errorf("Expected 15 total errors, got %d", hm.GetTotalErrors())
	}

	expectedAccuracy := (175.0 - 15.0) / 175.0 * 100.0
	if hm.GetOverallAccuracy() != expectedAccuracy {
		t.Errorf("Expected %.1f%% accuracy, got %.1f%%", expectedAccuracy, hm.GetOverallAccuracy())
	}
}

func TestHeatmapClear(t *testing.T) {
	hm := &Heatmap{
		Keys: map[string]*KeyStats{
			"a": {Key: "a", TotalHits: 100, ErrorCount: 10},
			"b": {Key: "b", TotalHits: 50, ErrorCount: 5},
		},
		path: filepath.Join(os.TempDir(), "test_heatmap.json"),
	}

	hm.Clear()

	if len(hm.Keys) != 0 {
		t.Errorf("Expected 0 keys after clear, got %d", len(hm.Keys))
	}

	if hm.GetTotalKeystrokes() != 0 {
		t.Error("Total keystrokes should be 0 after clear")
	}
}

func TestHeatmapSaveAndLoad(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "heatmap.json")

	// Create and populate heatmap
	hm1 := &Heatmap{
		Keys: map[string]*KeyStats{
			"a": {Key: "a", TotalHits: 10, ErrorCount: 2, LastUsed: time.Now()},
			"b": {Key: "b", TotalHits: 5, ErrorCount: 1, LastUsed: time.Now()},
		},
		path: path,
	}

	err := hm1.save()
	if err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	// Load into new heatmap
	hm2 := &Heatmap{
		Keys: make(map[string]*KeyStats),
		path: path,
	}
	hm2.load()

	if len(hm2.Keys) != 2 {
		t.Errorf("Expected 2 keys after load, got %d", len(hm2.Keys))
	}

	if hm2.Keys["a"].TotalHits != 10 {
		t.Errorf("Expected 'a' to have 10 hits after load, got %d", hm2.Keys["a"].TotalHits)
	}
}
