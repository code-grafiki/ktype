package storage

import (
	"fmt"
	"sort"
	"time"
)

// Statistics provides comprehensive typing analytics
type Statistics struct {
	lb *Leaderboard
}

// NewStatistics creates a new statistics analyzer
func NewStatistics(lb *Leaderboard) *Statistics {
	return &Statistics{lb: lb}
}

// StatsSummary provides a summary of typing statistics
type StatsSummary struct {
	TotalTests      int
	AverageWPM      float64
	AverageAccuracy float64
	BestWPM         int
	BestAccuracy    int
	RecentTests     int // Tests in last 7 days
}

// GetSummary returns overall statistics summary
func (s *Statistics) GetSummary() StatsSummary {
	if len(s.lb.Scores) == 0 {
		return StatsSummary{}
	}

	totalWPM := 0
	totalAccuracy := 0
	bestWPM := 0
	bestAccuracy := 0
	recentTests := 0
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)

	for _, score := range s.lb.Scores {
		totalWPM += score.WPM
		totalAccuracy += score.Accuracy

		if score.WPM > bestWPM {
			bestWPM = score.WPM
		}
		if score.Accuracy > bestAccuracy {
			bestAccuracy = score.Accuracy
		}

		if score.Date.After(sevenDaysAgo) {
			recentTests++
		}
	}

	count := len(s.lb.Scores)
	return StatsSummary{
		TotalTests:      count,
		AverageWPM:      float64(totalWPM) / float64(count),
		AverageAccuracy: float64(totalAccuracy) / float64(count),
		BestWPM:         bestWPM,
		BestAccuracy:    bestAccuracy,
		RecentTests:     recentTests,
	}
}

// ModeStats provides statistics for a specific mode
type ModeStats struct {
	Mode            string
	TestsCompleted  int
	AverageWPM      float64
	AverageAccuracy float64
	PersonalBest    int
	LastTestDate    time.Time
}

// GetModeStats returns statistics for each mode
func (s *Statistics) GetModeStats() []ModeStats {
	// Group scores by mode
	modeScores := make(map[string][]Score)
	for _, score := range s.lb.Scores {
		modeScores[score.Mode] = append(modeScores[score.Mode], score)
	}

	var stats []ModeStats
	for mode, scores := range modeScores {
		if len(scores) == 0 {
			continue
		}

		totalWPM := 0
		totalAccuracy := 0
		bestWPM := 0
		var lastDate time.Time

		for _, score := range scores {
			totalWPM += score.WPM
			totalAccuracy += score.Accuracy
			if score.WPM > bestWPM {
				bestWPM = score.WPM
			}
			if score.Date.After(lastDate) {
				lastDate = score.Date
			}
		}

		count := len(scores)
		stats = append(stats, ModeStats{
			Mode:            mode,
			TestsCompleted:  count,
			AverageWPM:      float64(totalWPM) / float64(count),
			AverageAccuracy: float64(totalAccuracy) / float64(count),
			PersonalBest:    bestWPM,
			LastTestDate:    lastDate,
		})
	}

	// Sort by mode name for consistent display
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Mode < stats[j].Mode
	})

	return stats
}

// TrendPoint represents a data point for trend analysis
type TrendPoint struct {
	Date     time.Time
	WPM      int
	Accuracy int
}

// GetWPMTrend returns WPM trend over time (last 30 tests)
func (s *Statistics) GetWPMTrend(limit int) []TrendPoint {
	if len(s.lb.Scores) == 0 {
		return []TrendPoint{}
	}

	// Get last N scores, sorted by date
	scores := make([]Score, len(s.lb.Scores))
	copy(scores, s.lb.Scores)

	// Sort by date ascending
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Date.Before(scores[j].Date)
	})

	// Take last 'limit' scores
	if len(scores) > limit {
		scores = scores[len(scores)-limit:]
	}

	var trend []TrendPoint
	for _, score := range scores {
		trend = append(trend, TrendPoint{
			Date:     score.Date,
			WPM:      score.WPM,
			Accuracy: score.Accuracy,
		})
	}

	return trend
}

// PerformanceRange categorizes performance levels
type PerformanceRange struct {
	Label string
	Min   int
	Max   int
	Count int
	Color string
}

// GetWPMDistribution returns distribution of WPM ranges
func (s *Statistics) GetWPMDistribution() []PerformanceRange {
	ranges := []PerformanceRange{
		{Label: "Beginner", Min: 0, Max: 30, Color: "#ca4754"},      // Red
		{Label: "Intermediate", Min: 30, Max: 50, Color: "#e2b714"}, // Yellow
		{Label: "Advanced", Min: 50, Max: 80, Color: "#5eacd3"},     // Blue
		{Label: "Expert", Min: 80, Max: 120, Color: "#98c379"},      // Green
		{Label: "Master", Min: 120, Max: 999, Color: "#c678dd"},     // Purple
	}

	for i := range ranges {
		count := 0
		for _, score := range s.lb.Scores {
			if score.WPM >= ranges[i].Min && score.WPM < ranges[i].Max {
				count++
			}
		}
		ranges[i].Count = count
	}

	return ranges
}

// ConsistencyMetrics measures typing consistency
type ConsistencyMetrics struct {
	WPMStdDev         float64 // Standard deviation of WPM
	AccuracyStdDev    float64 // Standard deviation of accuracy
	ConsistencyRating string  // "Very Consistent", "Consistent", "Variable", "Unpredictable"
}

// GetConsistencyMetrics calculates consistency metrics
func (s *Statistics) GetConsistencyMetrics() ConsistencyMetrics {
	if len(s.lb.Scores) < 3 {
		return ConsistencyMetrics{ConsistencyRating: "N/A"}
	}

	// Calculate mean
	sumWPM := 0
	sumAccuracy := 0
	for _, score := range s.lb.Scores {
		sumWPM += score.WPM
		sumAccuracy += score.Accuracy
	}
	meanWPM := float64(sumWPM) / float64(len(s.lb.Scores))
	meanAccuracy := float64(sumAccuracy) / float64(len(s.lb.Scores))

	// Calculate variance
	varianceWPM := 0.0
	varianceAccuracy := 0.0
	for _, score := range s.lb.Scores {
		diffWPM := float64(score.WPM) - meanWPM
		varianceWPM += diffWPM * diffWPM
		diffAccuracy := float64(score.Accuracy) - meanAccuracy
		varianceAccuracy += diffAccuracy * diffAccuracy
	}

	// Standard deviation
	stdDevWPM := varianceWPM / float64(len(s.lb.Scores))
	stdDevAccuracy := varianceAccuracy / float64(len(s.lb.Scores))

	// Determine rating based on WPM std dev
	rating := "Unpredictable"
	if stdDevWPM < 5 {
		rating = "Very Consistent"
	} else if stdDevWPM < 10 {
		rating = "Consistent"
	} else if stdDevWPM < 20 {
		rating = "Variable"
	}

	return ConsistencyMetrics{
		WPMStdDev:         stdDevWPM,
		AccuracyStdDev:    stdDevAccuracy,
		ConsistencyRating: rating,
	}
}

// RecentPerformance shows performance in recent time periods
type RecentPerformance struct {
	Today           int // Tests today
	ThisWeek        int // Tests this week
	ThisMonth       int // Tests this month
	AverageWPMToday float64
}

// GetRecentPerformance returns recent activity metrics
func (s *Statistics) GetRecentPerformance() RecentPerformance {
	now := time.Now()
	today := now.Truncate(24 * time.Hour)
	weekAgo := now.AddDate(0, 0, -7)
	monthAgo := now.AddDate(0, -1, 0)

	todayCount := 0
	weekCount := 0
	monthCount := 0
	todayWPM := 0

	for _, score := range s.lb.Scores {
		if score.Date.After(today) {
			todayCount++
			todayWPM += score.WPM
		}
		if score.Date.After(weekAgo) {
			weekCount++
		}
		if score.Date.After(monthAgo) {
			monthCount++
		}
	}

	averageToday := 0.0
	if todayCount > 0 {
		averageToday = float64(todayWPM) / float64(todayCount)
	}

	return RecentPerformance{
		Today:           todayCount,
		ThisWeek:        weekCount,
		ThisMonth:       monthCount,
		AverageWPMToday: averageToday,
	}
}

// FormatDuration formats a duration in a human-readable way
func FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	return fmt.Sprintf("%dh", int(d.Hours()))
}
