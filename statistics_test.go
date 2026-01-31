package main

import (
	"testing"
	"time"
)

func TestNewStatistics(t *testing.T) {
	lb := &Leaderboard{
		Scores: []Score{},
	}
	stats := NewStatistics(lb)
	if stats == nil {
		t.Fatal("NewStatistics returned nil")
	}
	if stats.lb != lb {
		t.Error("Statistics should reference the provided leaderboard")
	}
}

func TestGetSummaryEmpty(t *testing.T) {
	lb := &Leaderboard{
		Scores: []Score{},
	}
	stats := NewStatistics(lb)
	summary := stats.GetSummary()

	if summary.TotalTests != 0 {
		t.Errorf("Expected 0 tests, got %d", summary.TotalTests)
	}
}

func TestGetSummary(t *testing.T) {
	now := time.Now()
	lb := &Leaderboard{
		Scores: []Score{
			{WPM: 60, Accuracy: 95, Mode: "time:30", Date: now},
			{WPM: 70, Accuracy: 92, Mode: "time:60", Date: now},
			{WPM: 65, Accuracy: 98, Mode: "words:25", Date: now},
		},
	}
	stats := NewStatistics(lb)
	summary := stats.GetSummary()

	if summary.TotalTests != 3 {
		t.Errorf("Expected 3 tests, got %d", summary.TotalTests)
	}

	expectedAvgWPM := 65.0
	if summary.AverageWPM != expectedAvgWPM {
		t.Errorf("Expected average WPM %.1f, got %.1f", expectedAvgWPM, summary.AverageWPM)
	}

	expectedAvgAcc := 95.0
	if summary.AverageAccuracy != expectedAvgAcc {
		t.Errorf("Expected average accuracy %.1f, got %.1f", expectedAvgAcc, summary.AverageAccuracy)
	}

	if summary.BestWPM != 70 {
		t.Errorf("Expected best WPM 70, got %d", summary.BestWPM)
	}

	if summary.BestAccuracy != 98 {
		t.Errorf("Expected best accuracy 98, got %d", summary.BestAccuracy)
	}
}

func TestGetModeStats(t *testing.T) {
	now := time.Now()
	lb := &Leaderboard{
		Scores: []Score{
			{WPM: 60, Accuracy: 95, Mode: "time:30", Date: now},
			{WPM: 70, Accuracy: 92, Mode: "time:30", Date: now},
			{WPM: 65, Accuracy: 98, Mode: "words:25", Date: now},
		},
	}
	stats := NewStatistics(lb)
	modeStats := stats.GetModeStats()

	if len(modeStats) != 2 {
		t.Errorf("Expected 2 modes, got %d", len(modeStats))
	}

	// Check that modes are sorted
	if modeStats[0].Mode > modeStats[1].Mode {
		t.Error("Mode stats should be sorted by mode name")
	}

	// Find time:30 stats
	var time30Stats *ModeStats
	for i := range modeStats {
		if modeStats[i].Mode == "time:30" {
			time30Stats = &modeStats[i]
			break
		}
	}

	if time30Stats == nil {
		t.Fatal("time:30 mode stats not found")
	}

	if time30Stats.TestsCompleted != 2 {
		t.Errorf("Expected 2 tests for time:30, got %d", time30Stats.TestsCompleted)
	}

	expectedAvg := 65.0
	if time30Stats.AverageWPM != expectedAvg {
		t.Errorf("Expected average WPM %.1f for time:30, got %.1f", expectedAvg, time30Stats.AverageWPM)
	}

	if time30Stats.PersonalBest != 70 {
		t.Errorf("Expected personal best 70 for time:30, got %d", time30Stats.PersonalBest)
	}
}

func TestGetWPMTrend(t *testing.T) {
	now := time.Now()
	lb := &Leaderboard{
		Scores: []Score{
			{WPM: 50, Accuracy: 90, Mode: "time:30", Date: now.Add(-2 * time.Hour)},
			{WPM: 60, Accuracy: 92, Mode: "time:30", Date: now.Add(-1 * time.Hour)},
			{WPM: 70, Accuracy: 95, Mode: "time:30", Date: now},
		},
	}
	stats := NewStatistics(lb)
	trend := stats.GetWPMTrend(10)

	if len(trend) != 3 {
		t.Errorf("Expected 3 trend points, got %d", len(trend))
	}

	// Check that trend is sorted by date ascending
	for i := 1; i < len(trend); i++ {
		if trend[i].Date.Before(trend[i-1].Date) {
			t.Error("Trend should be sorted by date ascending")
		}
	}

	// Check values
	if trend[0].WPM != 50 || trend[1].WPM != 60 || trend[2].WPM != 70 {
		t.Error("Trend WPM values don't match expected")
	}
}

func TestGetWPMTrendLimit(t *testing.T) {
	now := time.Now()
	lb := &Leaderboard{
		Scores: []Score{},
	}
	// Add 20 scores
	for i := 0; i < 20; i++ {
		lb.Scores = append(lb.Scores, Score{
			WPM:      50 + i,
			Accuracy: 90,
			Mode:     "time:30",
			Date:     now.Add(time.Duration(i) * time.Hour),
		})
	}
	stats := NewStatistics(lb)
	trend := stats.GetWPMTrend(5)

	if len(trend) != 5 {
		t.Errorf("Expected 5 trend points (limited), got %d", len(trend))
	}
}

func TestGetWPMDistribution(t *testing.T) {
	lb := &Leaderboard{
		Scores: []Score{
			{WPM: 20, Accuracy: 90, Mode: "time:30", Date: time.Now()},  // Beginner
			{WPM: 35, Accuracy: 90, Mode: "time:30", Date: time.Now()},  // Intermediate
			{WPM: 60, Accuracy: 90, Mode: "time:30", Date: time.Now()},  // Advanced
			{WPM: 90, Accuracy: 90, Mode: "time:30", Date: time.Now()},  // Expert
			{WPM: 130, Accuracy: 90, Mode: "time:30", Date: time.Now()}, // Master
		},
	}
	stats := NewStatistics(lb)
	distribution := stats.GetWPMDistribution()

	if len(distribution) != 5 {
		t.Errorf("Expected 5 performance ranges, got %d", len(distribution))
	}

	// Check that each range has count 1
	for _, r := range distribution {
		if r.Count != 1 {
			t.Errorf("Expected count 1 for %s, got %d", r.Label, r.Count)
		}
	}
}

func TestGetConsistencyMetrics(t *testing.T) {
	now := time.Now()
	lb := &Leaderboard{
		Scores: []Score{
			{WPM: 60, Accuracy: 95, Mode: "time:30", Date: now},
			{WPM: 61, Accuracy: 94, Mode: "time:30", Date: now},
			{WPM: 59, Accuracy: 96, Mode: "time:30", Date: now},
			{WPM: 60, Accuracy: 95, Mode: "time:30", Date: now},
			{WPM: 60, Accuracy: 95, Mode: "time:30", Date: now},
		},
	}
	stats := NewStatistics(lb)
	consistency := stats.GetConsistencyMetrics()

	// The standard deviation should be very low (around 0.5 or less)
	// so we accept anything that rates as "Very Consistent"
	if consistency.ConsistencyRating != "Very Consistent" {
		t.Errorf("Expected 'Very Consistent' for consistent scores, got %s (std dev: %.2f)",
			consistency.ConsistencyRating, consistency.WPMStdDev)
	}

	// Test with more variable scores
	lb2 := &Leaderboard{
		Scores: []Score{
			{WPM: 30, Accuracy: 85, Mode: "time:30", Date: now},
			{WPM: 60, Accuracy: 95, Mode: "time:30", Date: now},
			{WPM: 90, Accuracy: 98, Mode: "time:30", Date: now},
		},
	}
	stats2 := NewStatistics(lb2)
	consistency2 := stats2.GetConsistencyMetrics()

	// Standard deviation should be around 600 for these spread out values
	if consistency2.ConsistencyRating == "Very Consistent" || consistency2.ConsistencyRating == "Consistent" {
		t.Errorf("Expected inconsistent rating for spread out scores, got %s (std dev: %.2f)",
			consistency2.ConsistencyRating, consistency2.WPMStdDev)
	}
}

func TestGetConsistencyMetricsEmpty(t *testing.T) {
	lb := &Leaderboard{
		Scores: []Score{},
	}
	stats := NewStatistics(lb)
	consistency := stats.GetConsistencyMetrics()

	if consistency.ConsistencyRating != "N/A" {
		t.Errorf("Expected 'N/A' for empty data, got %s", consistency.ConsistencyRating)
	}
}

func TestGetRecentPerformance(t *testing.T) {
	now := time.Now()
	lb := &Leaderboard{
		Scores: []Score{
			{WPM: 60, Accuracy: 95, Mode: "time:30", Date: now},                          // Today
			{WPM: 65, Accuracy: 96, Mode: "time:30", Date: now.Add(-24 * time.Hour)},     // Yesterday
			{WPM: 70, Accuracy: 97, Mode: "time:30", Date: now.Add(-6 * 24 * time.Hour)}, // Within week
			{WPM: 75, Accuracy: 98, Mode: "time:30", Date: now.AddDate(0, 0, -29)},       // Within month
			{WPM: 80, Accuracy: 99, Mode: "time:30", Date: now.AddDate(0, 0, -35)},       // Outside month
		},
	}
	stats := NewStatistics(lb)
	recent := stats.GetRecentPerformance()

	if recent.Today != 1 {
		t.Errorf("Expected 1 test today, got %d", recent.Today)
	}

	// This week should include today, yesterday, and 6 days ago
	if recent.ThisWeek != 3 {
		t.Errorf("Expected 3 tests this week, got %d", recent.ThisWeek)
	}

	// This month should include all except the one 35 days ago
	if recent.ThisMonth != 4 {
		t.Errorf("Expected 4 tests this month, got %d", recent.ThisMonth)
	}

	if recent.AverageWPMToday != 60.0 {
		t.Errorf("Expected 60.0 avg WPM today, got %.1f", recent.AverageWPMToday)
	}
}

func TestRenderBar(t *testing.T) {
	// Test basic bar rendering
	bar := renderBar(10, 20, 10)
	if len(bar) == 0 {
		t.Error("renderBar should return non-empty string")
	}

	// Test edge cases
	bar0 := renderBar(0, 100, 10)
	if len(bar0) == 0 {
		t.Error("renderBar should handle zero value")
	}

	barMax := renderBar(100, 100, 10)
	if len(barMax) == 0 {
		t.Error("renderBar should handle max value")
	}
}
