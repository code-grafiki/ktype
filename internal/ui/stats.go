package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"ktype/internal/storage"
)

// RenderStats renders the statistics dashboard
func RenderStats(stats *storage.Statistics, width, height int, wantToQuit bool) string {
	var s strings.Builder

	title := titleStyle.Render("statistics")
	s.WriteString(title)
	s.WriteString("\n\n")

	// Summary statistics
	summary := stats.GetSummary()
	if summary.TotalTests == 0 {
		s.WriteString(subtleStyle.Render("no data yet - complete some tests to see statistics"))
		s.WriteString("\n\n")
	} else {
		// Overall stats
		s.WriteString(subtleStyle.Render("overall:"))
		s.WriteString("\n")
		s.WriteString(fmt.Sprintf("   %s %s\n",
			wpmStyle.Render(fmt.Sprintf("%.1f", summary.AverageWPM)),
			subtleStyle.Render("avg wpm")))
		s.WriteString(fmt.Sprintf("   %s %s\n",
			accuracyStyle.Render(fmt.Sprintf("%.1f%%", summary.AverageAccuracy)),
			subtleStyle.Render("avg accuracy")))
		s.WriteString(fmt.Sprintf("   %s %s\n",
			wpmStyle.Render(fmt.Sprintf("%d", summary.BestWPM)),
			subtleStyle.Render("best wpm")))
		s.WriteString(fmt.Sprintf("   %s %s\n",
			statsStyle.Render(fmt.Sprintf("%d", summary.TotalTests)),
			subtleStyle.Render("total tests")))
		s.WriteString(fmt.Sprintf("   %s %s\n",
			statsStyle.Render(fmt.Sprintf("%d", summary.RecentTests)),
			subtleStyle.Render("tests this week")))
		s.WriteString("\n")

		// Recent activity
		recent := stats.GetRecentPerformance()
		s.WriteString(subtleStyle.Render("recent activity:"))
		s.WriteString("\n")
		s.WriteString(fmt.Sprintf("   %s %s\n",
			statsStyle.Render(fmt.Sprintf("%d", recent.Today)),
			subtleStyle.Render("tests today")))
		s.WriteString(fmt.Sprintf("   %s %s\n",
			statsStyle.Render(fmt.Sprintf("%d", recent.ThisWeek)),
			subtleStyle.Render("tests this week")))
		s.WriteString(fmt.Sprintf("   %s %s\n",
			statsStyle.Render(fmt.Sprintf("%d", recent.ThisMonth)),
			subtleStyle.Render("tests this month")))
		if recent.Today > 0 {
			s.WriteString(fmt.Sprintf("   %s %s\n",
				wpmStyle.Render(fmt.Sprintf("%.1f", recent.AverageWPMToday)),
				subtleStyle.Render("avg wpm today")))
		}
		s.WriteString("\n")

		// Consistency metrics
		consistency := stats.GetConsistencyMetrics()
		s.WriteString(subtleStyle.Render("consistency:"))
		s.WriteString("\n")
		s.WriteString(fmt.Sprintf("   %s %s\n",
			statsStyle.Render(consistency.ConsistencyRating),
			subtleStyle.Render("rating")))
		s.WriteString("\n")

		// Mode statistics
		modeStats := stats.GetModeStats()
		if len(modeStats) > 0 {
			s.WriteString(subtleStyle.Render("by mode:"))
			s.WriteString("\n")
			for _, ms := range modeStats {
				modeLine := fmt.Sprintf("   %s: %s %s",
					subtleStyle.Render(ms.Mode),
					wpmStyle.Render(fmt.Sprintf("%.0f", ms.AverageWPM)),
					subtleStyle.Render("wpm avg"))
				s.WriteString(modeLine)
				s.WriteString("\n")
			}
			s.WriteString("\n")
		}

		// WPM Distribution
		distribution := stats.GetWPMDistribution()
		s.WriteString(subtleStyle.Render("performance distribution:"))
		s.WriteString("\n")
		for _, r := range distribution {
			if r.Count > 0 {
				bar := renderBar(r.Count, summary.TotalTests, 20)
				s.WriteString(fmt.Sprintf("   %s: %s %s\n",
					subtleStyle.Render(r.Label),
					bar,
					statsStyle.Render(fmt.Sprintf("%d", r.Count))))
			}
		}
	}

	s.WriteString("\n")
	var help string
	if wantToQuit {
		help = errorStyle.Render("press esc again to go back")
	} else {
		help = helpStyle.Render("esc to go back")
	}
	s.WriteString(help)

	content := containerStyle.Render(s.String())
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}

// RenderHeatmap renders the typing heatmap visualization
func RenderHeatmap(hm *storage.Heatmap, width, height int, wantToQuit bool) string {
	var s strings.Builder

	title := titleStyle.Render("typing heatmap")
	s.WriteString(title)
	s.WriteString("\n\n")

	if hm.GetTotalKeystrokes() == 0 {
		s.WriteString(subtleStyle.Render("no data yet - type some words to see your heatmap"))
		s.WriteString("\n\n")
	} else {
		// Overall stats
		s.WriteString(subtleStyle.Render("overall:"))
		s.WriteString("\n")
		s.WriteString(fmt.Sprintf("   %s %s\n",
			statsStyle.Render(fmt.Sprintf("%d", hm.GetTotalKeystrokes())),
			subtleStyle.Render("total keystrokes")))
		s.WriteString(fmt.Sprintf("   %s %s\n",
			statsStyle.Render(fmt.Sprintf("%d", hm.GetTotalErrors())),
			subtleStyle.Render("total errors")))
		s.WriteString(fmt.Sprintf("   %s %s\n",
			accuracyStyle.Render(fmt.Sprintf("%.1f%%", hm.GetOverallAccuracy())),
			subtleStyle.Render("accuracy")))
		s.WriteString("\n")

		// Top error keys
		topErrors := hm.GetTopErrors(5)
		if len(topErrors) > 0 {
			s.WriteString(subtleStyle.Render("keys with most errors:"))
			s.WriteString("\n")
			for _, stat := range topErrors {
				errorRate := stat.ErrorRate()
				level := storage.GetErrorHeatLevel(errorRate)
				color := storage.GetHeatColor(level)
				keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Bold(true)

				s.WriteString(fmt.Sprintf("   %s %s %s\n",
					keyStyle.Render(fmt.Sprintf("[%s]", stat.Key)),
					subtleStyle.Render(fmt.Sprintf("%.1f%% error rate (%d/%d)", errorRate, stat.ErrorCount, stat.TotalHits)),
					renderBar(int(errorRate), 100, 10)))
			}
			s.WriteString("\n")
		}

		// Most used keys
		mostUsed := hm.GetMostUsed(5)
		if len(mostUsed) > 0 {
			s.WriteString(subtleStyle.Render("most used keys:"))
			s.WriteString("\n")
			for _, stat := range mostUsed {
				s.WriteString(fmt.Sprintf("   %s %s %s\n",
					wpmStyle.Render(fmt.Sprintf("[%s]", stat.Key)),
					subtleStyle.Render(fmt.Sprintf("%d hits", stat.TotalHits)),
					renderBar(stat.TotalHits, mostUsed[0].TotalHits, 10)))
			}
			s.WriteString("\n")
		}

		// Keyboard heatmap visualization
		s.WriteString(subtleStyle.Render("keyboard layout (error heat):"))
		s.WriteString("\n\n")

		keyboard := hm.GetHeatmapData()

		// Numbers row
		s.WriteString(subtleStyle.Render("  numbers: "))
		for _, stat := range keyboard.Numbers {
			errorRate := stat.ErrorRate()
			level := storage.GetErrorHeatLevel(errorRate)
			color := storage.GetHeatColor(level)
			keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
			s.WriteString(keyStyle.Render(fmt.Sprintf(" %s ", stat.Key)))
		}
		s.WriteString("\n\n")

		// Top row (QWERTY)
		s.WriteString(subtleStyle.Render("  top:     "))
		for _, stat := range keyboard.TopRow {
			errorRate := stat.ErrorRate()
			level := storage.GetErrorHeatLevel(errorRate)
			color := storage.GetHeatColor(level)
			keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
			s.WriteString(keyStyle.Render(fmt.Sprintf(" %s ", stat.Key)))
		}
		s.WriteString("\n\n")

		// Home row
		s.WriteString(subtleStyle.Render("  home:    "))
		for _, stat := range keyboard.HomeRow {
			errorRate := stat.ErrorRate()
			level := storage.GetErrorHeatLevel(errorRate)
			color := storage.GetHeatColor(level)
			keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
			s.WriteString(keyStyle.Render(fmt.Sprintf(" %s ", stat.Key)))
		}
		s.WriteString("\n\n")

		// Bottom row
		s.WriteString(subtleStyle.Render("  bottom:  "))
		for _, stat := range keyboard.BottomRow {
			errorRate := stat.ErrorRate()
			level := storage.GetErrorHeatLevel(errorRate)
			color := storage.GetHeatColor(level)
			keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
			s.WriteString(keyStyle.Render(fmt.Sprintf(" %s ", stat.Key)))
		}
		s.WriteString("\n\n")

		// Legend
		s.WriteString(subtleStyle.Render("legend: "))
		legend := []struct {
			color string
			label string
		}{
			{"#646669", "no errors"},
			{"#98c379", "low (<5%)"},
			{"#e2b714", "medium (5-15%)"},
			{"#d19a66", "high (15-30%)"},
			{"#ca4754", "very high (>30%)"},
		}
		for _, l := range legend {
			keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(l.color))
			s.WriteString(keyStyle.Render("█") + subtleStyle.Render(" "+l.label+"  "))
		}
		s.WriteString("\n")
	}

	s.WriteString("\n")
	var help string
	if wantToQuit {
		help = errorStyle.Render("press esc again to go back")
	} else {
		help = helpStyle.Render("esc to go back • r to reset")
	}
	s.WriteString(help)

	content := containerStyle.Render(s.String())
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}

// RenderChallenges renders the daily challenges screen
func RenderChallenges(dc *storage.DailyChallenges, width, height int, wantToQuit bool) string {
	var s strings.Builder

	title := titleStyle.Render("daily challenges")
	s.WriteString(title)
	s.WriteString("\n\n")

	challenges := dc.GetTodaysChallenges()

	if len(challenges) == 0 {
		s.WriteString(subtleStyle.Render("no challenges available"))
	} else {
		completed := dc.GetCompletedCount()
		total := dc.GetTotalCount()

		s.WriteString(subtleStyle.Render(fmt.Sprintf("progress: %d/%d completed", completed, total)))
		s.WriteString("\n\n")

		for i, c := range challenges {
			status := ""
			if c.Completed {
				status = accuracyStyle.Render(" ✓ COMPLETED")
			} else {
				progress := float64(c.Progress) / float64(c.Target) * 100
				if progress > 0 {
					status = subtleStyle.Render(fmt.Sprintf(" (%.0f%%)", progress))
				}
			}

			s.WriteString(fmt.Sprintf("%s %s%s\n",
				wpmStyle.Render(fmt.Sprintf("%d.", i+1)),
				subtleStyle.Render(c.Title),
				status))

			s.WriteString(fmt.Sprintf("   %s\n", subtleStyle.Render(c.Description)))

			if !c.Completed {
				s.WriteString(fmt.Sprintf("   %s: %s/%d %s\n",
					subtleStyle.Render("progress"),
					statsStyle.Render(fmt.Sprintf("%d", c.Progress)),
					c.Target,
					subtleStyle.Render("("+c.Reward+")")))
			}

			s.WriteString("\n")
		}

		if dc.HasCompletedAll() {
			s.WriteString("\n")
			s.WriteString(accuracyStyle.Render("🎉 All challenges completed today!"))
			s.WriteString("\n")
		}
	}

	s.WriteString("\n")
	var help string
	if wantToQuit {
		help = errorStyle.Render("press esc again to go back")
	} else {
		help = helpStyle.Render("esc to go back")
	}
	s.WriteString(help)

	content := containerStyle.Render(s.String())
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}
