package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// MonkeyType-inspired color palette
var (
	colorBg      = lipgloss.Color("#323437")
	colorSubtle  = lipgloss.Color("#646669")
	colorText    = lipgloss.Color("#d1d0c5")
	colorError   = lipgloss.Color("#ca4754")
	colorAccent  = lipgloss.Color("#5eacd3")
	colorCorrect = lipgloss.Color("#98c379")
)

// Styles
var (
	baseStyle = lipgloss.NewStyle().
			Background(colorBg).
			Foreground(colorText)

	titleStyle = lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true).
			MarginBottom(1)

	subtleStyle = lipgloss.NewStyle().
			Foreground(colorSubtle)

	correctStyle = lipgloss.NewStyle().
			Foreground(colorSubtle)

	errorStyle = lipgloss.NewStyle().
			Foreground(colorError).
			Underline(true)

	currentStyle = lipgloss.NewStyle().
			Foreground(colorText).
			Bold(true)

	cursorStyle = lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(false)

	upcomingStyle = lipgloss.NewStyle().
			Foreground(colorSubtle)

	timerStyle = lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true)

	statsStyle = lipgloss.NewStyle().
			Foreground(colorText)

	wpmStyle = lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true)

	accuracyStyle = lipgloss.NewStyle().
			Foreground(colorCorrect)

	containerStyle = lipgloss.NewStyle().
			Padding(2, 4).
			Margin(1, 2).
			Height(20).
			Width(65)

	helpStyle = lipgloss.NewStyle().
			Foreground(colorSubtle).
			MarginTop(2)

	pbStyle = lipgloss.NewStyle().
		Foreground(colorCorrect).
		Bold(true)

	newPBStyle = lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true)
)

// UpdateAccentColor updates the accent color and all dependent styles
func UpdateAccentColor(hex string) {
	colorAccent = lipgloss.Color(hex)

	// Update all styles that use colorAccent
	titleStyle = lipgloss.NewStyle().
		Foreground(colorAccent).
		Bold(true).
		MarginBottom(1)

	cursorStyle = lipgloss.NewStyle().
		Foreground(colorAccent).
		Bold(false)

	timerStyle = lipgloss.NewStyle().
		Foreground(colorAccent).
		Bold(true)

	wpmStyle = lipgloss.NewStyle().
		Foreground(colorAccent).
		Bold(true)

	newPBStyle = lipgloss.NewStyle().
		Foreground(colorAccent).
		Bold(true)
}

// RenderMainMenu renders the main menu with quick start options
func RenderMainMenu(lb *Leaderboard, width, height int, wantToQuit bool, difficulty Difficulty, complexity WordComplexity) string {
	var s strings.Builder

	title := titleStyle.Render("ktype")
	s.WriteString(title)
	s.WriteString("\n\n")

	// Quick Start Presets
	s.WriteString(subtleStyle.Render("quick start:"))
	s.WriteString("\n\n")

	// Get PBs for presets
	pb30s := lb.GetPB("time:30")
	pb50w := lb.GetPB("words:50")
	pbZen := lb.GetPB("zen")

	formatPB := func(s *Score) string {
		if s == nil {
			return subtleStyle.Render(" (no PB yet)")
		}
		return pbStyle.Render(fmt.Sprintf(" (PB: %d wpm | %d%%)", s.WPM, s.Accuracy))
	}

	options := []string{
		wpmStyle.Render("1") + subtleStyle.Render(" → 30s timed") + formatPB(pb30s),
		wpmStyle.Render("2") + subtleStyle.Render(" → 50 words") + formatPB(pb50w),
		wpmStyle.Render("3") + subtleStyle.Render(" → zen mode") + formatPB(pbZen),
	}

	for _, opt := range options {
		s.WriteString("   " + opt + "\n")
	}

	s.WriteString("\n")
	// More Modes
	s.WriteString(subtleStyle.Render("more modes:"))
	s.WriteString("\n\n")

	moreModes := []string{
		wpmStyle.Render("t") + subtleStyle.Render(" → timed modes selection"),
		wpmStyle.Render("w") + subtleStyle.Render(" → words modes selection"),
	}

	for _, opt := range moreModes {
		s.WriteString("   " + opt + "\n")
	}

	// Difficulty
	s.WriteString("\n")
	s.WriteString(subtleStyle.Render("difficulty: ") + wpmStyle.Render(difficulty.String()))
	s.WriteString("\n")
	s.WriteString("   " + wpmStyle.Render("d") + subtleStyle.Render(" → change difficulty\n"))

	// Complexity
	s.WriteString("\n")
	s.WriteString(subtleStyle.Render("content: ") + wpmStyle.Render(complexity.String()))
	s.WriteString("\n")
	s.WriteString("   " + wpmStyle.Render("c") + subtleStyle.Render(" → change content (punctuation/numbers)\n"))

	// Statistics
	s.WriteString("\n")
	s.WriteString(subtleStyle.Render("statistics:"))
	s.WriteString("\n")
	s.WriteString("   " + wpmStyle.Render("s") + subtleStyle.Render(" → view statistics\n"))

	// Heatmap
	s.WriteString("\n")
	s.WriteString(subtleStyle.Render("analysis:"))
	s.WriteString("\n")
	s.WriteString("   " + wpmStyle.Render("h") + subtleStyle.Render(" → typing heatmap\n"))

	// Custom Word Lists
	s.WriteString("\n")
	s.WriteString(subtleStyle.Render("word lists:"))
	s.WriteString("\n")
	s.WriteString("   " + wpmStyle.Render("l") + subtleStyle.Render(" → custom word lists\n"))

	// Daily Challenges
	s.WriteString("\n")
	s.WriteString(subtleStyle.Render("daily challenges:"))
	s.WriteString("\n")
	s.WriteString("   " + wpmStyle.Render("v") + subtleStyle.Render(" → view challenges\n"))

	// Settings
	s.WriteString("\n")
	s.WriteString(subtleStyle.Render("settings:"))
	s.WriteString("\n")
	s.WriteString("   " + wpmStyle.Render(",") + subtleStyle.Render(" → appearance settings\n"))

	// Exit prompt
	s.WriteString("\n")
	var help string
	if wantToQuit {
		help = errorStyle.Render("press esc again to quit")
	} else {
		help = helpStyle.Render("esc to quit")
	}
	s.WriteString(help)

	content := containerStyle.Render(s.String())
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}

// RenderTimeSelect renders the time duration selection screen with PBs
func RenderTimeSelect(lb *Leaderboard, width, height int, wantToQuit bool) string {
	var s strings.Builder

	title := titleStyle.Render("timed mode")
	s.WriteString(title)
	s.WriteString("\n\n")

	s.WriteString(subtleStyle.Render("select duration:"))
	s.WriteString("\n\n")

	durations := []struct {
		key   string
		label string
		mode  string
	}{
		{"1", "15s", "time:15"},
		{"2", "30s", "time:30"},
		{"3", "60s", "time:60"},
	}

	for _, d := range durations {
		pb := lb.GetPB(d.mode)
		pbText := ""
		if pb != nil {
			pbText = pbStyle.Render(fmt.Sprintf(" (PB: %d|%d%%)", pb.WPM, pb.Accuracy))
		}
		s.WriteString(fmt.Sprintf("   %s %s%s\n", wpmStyle.Render(d.key), subtleStyle.Render("→ "+d.label), pbText))
	}
	s.WriteString("   " + wpmStyle.Render("c") + subtleStyle.Render(" → custom\n"))

	s.WriteString("\n")
	var help string
	if wantToQuit {
		help = errorStyle.Render("press esc again to quit")
	} else {
		help = helpStyle.Render("esc: back")
	}
	s.WriteString(help)

	content := containerStyle.Render(s.String())
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}

// RenderWordsSelect renders the word count selection screen with PBs
func RenderWordsSelect(lb *Leaderboard, width, height int, wantToQuit bool) string {
	var s strings.Builder

	title := titleStyle.Render("words mode")
	s.WriteString(title)
	s.WriteString("\n\n")

	s.WriteString(subtleStyle.Render("select word count:"))
	s.WriteString("\n\n")

	counts := []struct {
		key   string
		label string
		mode  string
	}{
		{"1", "10 words", "words:10"},
		{"2", "25 words", "words:25"},
		{"3", "50 words", "words:50"},
		{"4", "100 words", "words:100"},
	}

	for _, c := range counts {
		pb := lb.GetPB(c.mode)
		pbText := ""
		if pb != nil {
			pbText = pbStyle.Render(fmt.Sprintf(" (PB: %d|%d%%)", pb.WPM, pb.Accuracy))
		}
		s.WriteString(fmt.Sprintf("   %s %s%s\n", wpmStyle.Render(c.key), subtleStyle.Render("→ "+c.label), pbText))
	}
	s.WriteString("   " + wpmStyle.Render("c") + subtleStyle.Render(" → custom\n"))

	s.WriteString("\n")
	var help string
	if wantToQuit {
		help = errorStyle.Render("press esc again to quit")
	} else {
		help = helpStyle.Render("esc: back")
	}
	s.WriteString(help)

	content := containerStyle.Render(s.String())
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}

// RenderDifficultySelect renders the difficulty selection screen
func RenderDifficultySelect(currentDifficulty Difficulty, width, height int, wantToQuit bool) string {
	var s strings.Builder

	title := titleStyle.Render("difficulty")
	s.WriteString(title)
	s.WriteString("\n\n")

	s.WriteString(subtleStyle.Render("select difficulty:"))
	s.WriteString("\n\n")

	options := []struct {
		key        string
		label      string
		difficulty Difficulty
		desc       string
	}{
		{"1", "easy", DifficultyEasy, "2-4 letter words"},
		{"2", "medium", DifficultyMedium, "5-7 letter words"},
		{"3", "hard", DifficultyHard, "8+ letter words"},
	}

	for _, opt := range options {
		keyStyle := wpmStyle
		labelStyle := subtleStyle
		if opt.difficulty == currentDifficulty {
			keyStyle = lipgloss.NewStyle().Foreground(colorAccent).Bold(true)
			labelStyle = lipgloss.NewStyle().Foreground(colorText)
		}
		s.WriteString(fmt.Sprintf("   %s %s (%s)\n", keyStyle.Render(opt.key), labelStyle.Render("→ "+opt.label), opt.desc))
	}

	s.WriteString("\n")
	var help string
	if wantToQuit {
		help = errorStyle.Render("press esc again to quit")
	} else {
		help = helpStyle.Render("esc: back")
	}
	s.WriteString(help)

	content := containerStyle.Render(s.String())
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}

// RenderComplexitySelect renders the complexity selection screen
func RenderComplexitySelect(currentComplexity WordComplexity, width, height int, wantToQuit bool) string {
	var s strings.Builder

	title := titleStyle.Render("complexity")
	s.WriteString(title)
	s.WriteString("\n\n")

	s.WriteString(subtleStyle.Render("select complexity:"))
	s.WriteString("\n\n")

	options := []struct {
		key        string
		label      string
		complexity WordComplexity
		desc       string
	}{
		{"1", "normal", ComplexityNormal, "letters only"},
		{"2", "punctuation", ComplexityPunctuation, "letters + punctuation"},
		{"3", "numbers", ComplexityNumbers, "letters + numbers"},
		{"4", "full", ComplexityFull, "letters + punctuation + numbers"},
	}

	for _, opt := range options {
		keyStyle := wpmStyle
		labelStyle := subtleStyle
		if opt.complexity == currentComplexity {
			keyStyle = lipgloss.NewStyle().Foreground(colorAccent).Bold(true)
			labelStyle = lipgloss.NewStyle().Foreground(colorText)
		}
		s.WriteString(fmt.Sprintf("   %s %s (%s)\n", keyStyle.Render(opt.key), labelStyle.Render("→ "+opt.label), opt.desc))
	}

	s.WriteString("\n")
	var help string
	if wantToQuit {
		help = errorStyle.Render("press esc again to quit")
	} else {
		help = helpStyle.Render("esc: back")
	}
	s.WriteString(help)

	content := containerStyle.Render(s.String())
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}

// RenderCustomInput renders the custom input screen
func RenderCustomInput(input, mode string, width, height int, cursorChar string) string {
	var s strings.Builder

	var prompt string
	if mode == "time" {
		prompt = "enter duration (seconds)"
	} else {
		prompt = "enter word count"
	}

	title := titleStyle.Render(prompt)
	s.WriteString(title)
	s.WriteString("\n\n")

	// Show input with cursor
	inputDisplay := wpmStyle.Render(input) + cursorStyle.Render(cursorChar)
	s.WriteString("        " + inputDisplay)
	s.WriteString("\n\n")

	s.WriteString(helpStyle.Render("enter to confirm • esc to go back"))

	content := containerStyle.Render(s.String())
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}

// RenderGame renders the main game screen
func RenderGame(g *Game, width, height int, wantToQuit bool, cursorChar string) string {
	var s strings.Builder

	// Inside the container, we have Width(65) and Padding(2, 4).
	// So internal width is 65 - 4*2 = 57.
	internalWidth := 57

	// Build the words display (3 lines)
	// We use a slightly smaller width for the words themselves to ensure they fit well
	wordsLines := buildWordsLines(g, internalWidth-2, 3, cursorChar)

	s.WriteString("\n")
	for _, line := range wordsLines {
		s.WriteString(lipgloss.PlaceHorizontal(internalWidth, lipgloss.Center, line))
		s.WriteString("\n")
	}
	s.WriteString("\n")

	// Progress (timer or word count)
	progress := timerStyle.Render(g.Progress())
	s.WriteString(lipgloss.PlaceHorizontal(internalWidth, lipgloss.Center, progress))
	s.WriteString("\n\n")

	// Live stats
	wpm := wpmStyle.Render(fmt.Sprintf("%d wpm", g.WPM()))
	accuracy := accuracyStyle.Render(fmt.Sprintf("%d%%", g.Accuracy()))
	liveErrors := len(g.CurrentErrors)
	liveErrorDisplay := ""
	if liveErrors > 0 {
		liveErrorDisplay = errorStyle.Render(fmt.Sprintf("%d", liveErrors))
	} else {
		liveErrorDisplay = subtleStyle.Render("0")
	}
	stats := statsStyle.Render(wpm + subtleStyle.Render("  •  ") + accuracy + subtleStyle.Render("  •  ") + liveErrorDisplay + subtleStyle.Render(" errors"))
	s.WriteString(lipgloss.PlaceHorizontal(internalWidth, lipgloss.Center, stats))

	// Error statistics section (subtle)
	totalErrors, topErrors := GetErrorStats(g)
	if totalErrors > 0 {
		s.WriteString("\n")
		var errorStats strings.Builder
		errorStats.WriteString(errorDetailStyle.Render("errors: "))
		errorStats.WriteString(subtleStyle.Render(fmt.Sprintf("%d total", totalErrors)))
		if len(topErrors) > 0 {
			errorStats.WriteString(errorDetailStyle.Render(" • "))
			errorStats.WriteString(subtleStyle.Render("top: "))
			for i, err := range topErrors {
				if i > 0 {
					errorStats.WriteString(subtleStyle.Render(", "))
				}
				errorStats.WriteString(actualCharStyle.Render(string(err.Char)))
				errorStats.WriteString(subtleStyle.Render(fmt.Sprintf("(%d)", err.Count)))
			}
		}
		s.WriteString(lipgloss.PlaceHorizontal(internalWidth, lipgloss.Center, errorStats.String()))
	}

	// Help text or quit confirmation
	s.WriteString("\n\n")
	var help string
	if wantToQuit {
		help = errorStyle.Render("press esc again to quit")
	} else {
		help = helpStyle.Render("tab: restart • esc: abort")
	}
	s.WriteString(lipgloss.PlaceHorizontal(internalWidth, lipgloss.Center, help))

	content := containerStyle.Render(s.String())
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}

// RenderFinished renders the end screen
func RenderFinished(g *Game, width, height int, isPB bool, wantToQuit bool) string {
	var s strings.Builder

	if isPB {
		title := newPBStyle.Render("🎉  new personal best!")
		s.WriteString(title)
	} else {
		title := titleStyle.Render("test complete!")
		s.WriteString(title)
	}
	s.WriteString("\n\n")

	// Final WPM - big and prominent
	wpmValue := lipgloss.NewStyle().
		Foreground(colorAccent).
		Bold(true).
		Render(fmt.Sprintf("%d", g.WPM()))

	s.WriteString(wpmValue)
	s.WriteString(subtleStyle.Render(" wpm"))
	s.WriteString("\n\n")

	// Stats
	stats := []string{
		subtleStyle.Render("accuracy: ") + accuracyStyle.Render(fmt.Sprintf("%d%%", g.Accuracy())),
		subtleStyle.Render("raw wpm: ") + statsStyle.Render(fmt.Sprintf("%d", g.RawWPM())),
		subtleStyle.Render("correct: ") + statsStyle.Render(fmt.Sprintf("%d/%d words", g.CorrectWordsCount(), len(g.TypedWords))),
		subtleStyle.Render("mode: ") + statsStyle.Render(g.ModeString()),
	}

	for _, stat := range stats {
		s.WriteString(stat + "\n")
	}

	s.WriteString("\n")
	var help string
	if wantToQuit {
		help = errorStyle.Render("press esc again to quit")
	} else {
		help = helpStyle.Render("tab to restart • esc to quit")
	}
	s.WriteString(help)

	content := containerStyle.Render(s.String())
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}

// buildWordsLines builds 3 lines of words with proper scrolling
func buildWordsLines(g *Game, maxWidth int, numLines int, cursorChar string) []string {
	type lineInfo struct {
		startIdx int
		endIdx   int
	}

	lines := []lineInfo{}
	currentLineStart := 0
	currentLineWidth := 0

	for i := 0; i < len(g.Words); i++ {
		wordWidth := len(g.Words[i])
		spaceNeeded := wordWidth
		if currentLineWidth > 0 {
			spaceNeeded += 1
		}

		if currentLineWidth+spaceNeeded > maxWidth && currentLineWidth > 0 {
			lines = append(lines, lineInfo{
				startIdx: currentLineStart,
				endIdx:   i,
			})
			currentLineStart = i
			currentLineWidth = wordWidth
		} else {
			currentLineWidth += spaceNeeded
		}
	}
	if currentLineStart < len(g.Words) {
		lines = append(lines, lineInfo{
			startIdx: currentLineStart,
			endIdx:   len(g.Words),
		})
	}

	currentLine := 0
	for i, line := range lines {
		if g.WordIndex >= line.startIdx && g.WordIndex < line.endIdx {
			currentLine = i
			break
		}
	}

	result := make([]string, numLines)

	for lineNum := 0; lineNum < numLines; lineNum++ {
		lineIdx := currentLine + lineNum
		if lineIdx >= len(lines) {
			result[lineNum] = ""
			continue
		}

		line := lines[lineIdx]
		var parts []string
		var rawWords []string

		for wordIdx := line.startIdx; wordIdx < line.endIdx; wordIdx++ {
			if wordIdx >= len(g.Words) {
				break
			}

			word := g.Words[wordIdx]
			rawWords = append(rawWords, word)

			if wordIdx < g.WordIndex {
				if wordIdx < len(g.Correct) && g.Correct[wordIdx] {
					parts = append(parts, correctStyle.Render(word))
				} else {
					parts = append(parts, errorStyle.Render(word))
				}
			} else if wordIdx == g.WordIndex {
				correct, errors, remaining := g.CurrentWordState()
				cursor := cursorStyle.Render(cursorChar)
				currentWord := correctStyle.Render(correct) +
					errorStyle.Render(errors) +
					cursor +
					currentStyle.Render(remaining)
				parts = append(parts, currentWord)
			} else {
				parts = append(parts, upcomingStyle.Render(word))
			}
		}

		isLastLine := lineIdx == len(lines)-1
		if !isLastLine && len(parts) > 1 {
			result[lineNum] = justifyLine(parts, rawWords, maxWidth)
		} else {
			result[lineNum] = strings.Join(parts, " ")
		}
	}

	return result
}

// justifyLine distributes extra spaces between words to fill maxWidth
func justifyLine(styledParts []string, rawWords []string, maxWidth int) string {
	if len(styledParts) <= 1 {
		return strings.Join(styledParts, " ")
	}

	totalWordLen := 0
	for _, w := range rawWords {
		totalWordLen += len(w)
	}
	totalWordLen += 1 // cursor

	gaps := len(styledParts) - 1
	totalSpaces := maxWidth - totalWordLen
	if totalSpaces < gaps {
		totalSpaces = gaps
	}

	baseSpaces := totalSpaces / gaps
	extraSpaces := totalSpaces % gaps

	var result strings.Builder
	for i, part := range styledParts {
		result.WriteString(part)
		if i < gaps {
			spaces := baseSpaces
			if i < extraSpaces {
				spaces++
			}
			result.WriteString(strings.Repeat(" ", spaces))
		}
	}

	return result.String()
}

// RenderStats renders the statistics dashboard
func RenderStats(stats *Statistics, width, height int, wantToQuit bool) string {
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

// renderBar creates a simple ASCII bar chart
func renderBar(value, max, width int) string {
	if max == 0 {
		return strings.Repeat("░", width)
	}
	filled := int(float64(value) / float64(max) * float64(width))
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#5eacd3")).
		Render(strings.Repeat("█", filled)) +
		lipgloss.NewStyle().Foreground(lipgloss.Color("#646669")).
			Render(strings.Repeat("░", width-filled))
}

// RenderHeatmap renders the typing heatmap visualization
func RenderHeatmap(hm *Heatmap, width, height int, wantToQuit bool) string {
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
				level := GetErrorHeatLevel(errorRate)
				color := GetHeatColor(level)
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
			level := GetErrorHeatLevel(errorRate)
			color := GetHeatColor(level)
			keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
			s.WriteString(keyStyle.Render(fmt.Sprintf(" %s ", stat.Key)))
		}
		s.WriteString("\n\n")

		// Top row (QWERTY)
		s.WriteString(subtleStyle.Render("  top:     "))
		for _, stat := range keyboard.TopRow {
			errorRate := stat.ErrorRate()
			level := GetErrorHeatLevel(errorRate)
			color := GetHeatColor(level)
			keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
			s.WriteString(keyStyle.Render(fmt.Sprintf(" %s ", stat.Key)))
		}
		s.WriteString("\n\n")

		// Home row
		s.WriteString(subtleStyle.Render("  home:    "))
		for _, stat := range keyboard.HomeRow {
			errorRate := stat.ErrorRate()
			level := GetErrorHeatLevel(errorRate)
			color := GetHeatColor(level)
			keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
			s.WriteString(keyStyle.Render(fmt.Sprintf(" %s ", stat.Key)))
		}
		s.WriteString("\n\n")

		// Bottom row
		s.WriteString(subtleStyle.Render("  bottom:  "))
		for _, stat := range keyboard.BottomRow {
			errorRate := stat.ErrorRate()
			level := GetErrorHeatLevel(errorRate)
			color := GetHeatColor(level)
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

// RenderCustomWordList renders the custom word list management screen
func RenderCustomWordList(wm *WordListManager, currentList string, width, height int, wantToQuit bool) string {
	var s strings.Builder

	title := titleStyle.Render("custom word lists")
	s.WriteString(title)
	s.WriteString("\n\n")

	if wm.Count() == 0 {
		s.WriteString(subtleStyle.Render("no custom word lists yet"))
		s.WriteString("\n\n")
		s.WriteString(subtleStyle.Render("word lists are stored in:"))
		s.WriteString("\n")
		s.WriteString(statsStyle.Render("~/.config/ktype/wordlists.json"))
		s.WriteString("\n\n")
	} else {
		s.WriteString(subtleStyle.Render("available word lists:"))
		s.WriteString("\n\n")

		for i, list := range wm.Lists {
			key := fmt.Sprintf("%d", i+1)
			selected := ""
			if list.Name == currentList {
				selected = accuracyStyle.Render(" ✓")
			}
			s.WriteString(fmt.Sprintf("   %s %s%s (%d words)\n",
				wpmStyle.Render(key),
				subtleStyle.Render(list.Name),
				selected,
				len(list.Words)))
			if list.Description != "" {
				s.WriteString(fmt.Sprintf("      %s\n", subtleStyle.Render(list.Description)))
			}
		}

		s.WriteString("\n")
		if currentList != "" {
			s.WriteString(accuracyStyle.Render(fmt.Sprintf("selected: %s", currentList)))
			s.WriteString("\n\n")
		}

		s.WriteString(subtleStyle.Render("actions:"))
		s.WriteString("\n")
		s.WriteString("   " + wpmStyle.Render("1-9") + subtleStyle.Render(" → select word list\n"))
		s.WriteString("   " + wpmStyle.Render("d") + subtleStyle.Render(" → delete selected\n"))
		s.WriteString("   " + wpmStyle.Render("c") + subtleStyle.Render(" → clear selection\n"))
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

// Enhanced error visualization colors
var (
	colorExpected = lipgloss.Color("#e2b714") // Yellow for expected char
	colorActual   = lipgloss.Color("#ca4754") // Red for actual typed char
	colorErrorDim = lipgloss.Color("#7c4c4c") // Dimmed red background
)

// Enhanced error styles
var (
	// expectedCharStyle shows what character should have been typed
	expectedCharStyle = lipgloss.NewStyle().
				Foreground(colorExpected).
				Bold(true)

	// actualCharStyle shows what was actually typed (wrong)
	actualCharStyle = lipgloss.NewStyle().
			Foreground(colorActual).
			Bold(true).
			Underline(true)

	// errorPositionStyle highlights the position of an error
	errorPositionStyle = lipgloss.NewStyle().
				Background(colorErrorDim).
				Foreground(colorText)

	// errorDetailStyle for showing error details panel
	errorDetailStyle = lipgloss.NewStyle().
				Foreground(colorSubtle).
				Italic(true)
)

// ErrorTypeString returns a human-readable error type description
func ErrorTypeString(et ErrorType) string {
	switch et {
	case ErrorWrongChar:
		return "wrong character"
	case ErrorExtraChar:
		return "extra character"
	case ErrorMissingChar:
		return "missing character"
	case ErrorTransposition:
		return "transposed"
	default:
		return "unknown"
	}
}

// FormatErrorDetails formats a typing error for display
func FormatErrorDetails(err TypingError) string {
	var b strings.Builder

	b.WriteString(errorDetailStyle.Render(fmt.Sprintf("[%s] ", ErrorTypeString(err.ErrorType))))

	if err.ErrorType == ErrorWrongChar {
		b.WriteString(fmt.Sprintf("expected %s but typed %s",
			expectedCharStyle.Render(string(err.ExpectedChar)),
			actualCharStyle.Render(string(err.TypedChar))))
	} else if err.ErrorType == ErrorExtraChar {
		b.WriteString(fmt.Sprintf("extra %s (word only has %d chars)",
			actualCharStyle.Render(string(err.TypedChar)),
			err.Position))
	}

	return b.String()
}

// RenderCurrentWordWithErrors renders the current word with detailed error highlighting
func RenderCurrentWordWithErrors(game *Game) string {
	if game.WordIndex >= len(game.Words) {
		return ""
	}

	targetWord := game.Words[game.WordIndex]
	input := game.CurrentInput

	if len(input) == 0 {
		return currentStyle.Render(targetWord)
	}

	var parts []string

	for i := 0; i < len(input) || i < len(targetWord); i++ {
		if i < len(input) && i < len(targetWord) {
			// Both input and target have character at this position
			if input[i] == targetWord[i] {
				// Correct character
				parts = append(parts, correctStyle.Render(string(input[i])))
			} else {
				// Wrong character - show both expected and actual
				parts = append(parts,
					actualCharStyle.Render(string(input[i]))+
						errorDetailStyle.Render("/")+
						expectedCharStyle.Render(string(targetWord[i])))
			}
		} else if i < len(input) {
			// Extra character beyond word length
			parts = append(parts, actualCharStyle.Render(string(input[i])))
		} else {
			// Remaining characters in target word
			parts = append(parts, currentStyle.Render(string(targetWord[i])))
		}
	}

	return strings.Join(parts, "")
}

// GetErrorStats returns statistics about errors in the current game
func GetErrorStats(game *Game) (totalErrors int, topErrors []struct {
	Char  rune
	Count int
}) {
	totalErrors = len(game.Errors)

	// Count errors by character
	charCounts := make(map[rune]int)
	for _, err := range game.Errors {
		charCounts[err.TypedChar]++
	}

	// Convert to slice for sorting
	type charCount struct {
		Char  rune
		Count int
	}

	var counts []charCount
	for char, count := range charCounts {
		counts = append(counts, charCount{Char: char, Count: count})
	}

	// Sort by count descending
	sort.Slice(counts, func(i, j int) bool {
		return counts[i].Count > counts[j].Count
	})

	// Return top 3
	if len(counts) > 3 {
		counts = counts[:3]
	}

	// Convert back to the expected type
	var result []struct {
		Char  rune
		Count int
	}
	for _, c := range counts {
		result = append(result, struct {
			Char  rune
			Count int
		}{Char: c.Char, Count: c.Count})
	}

	return totalErrors, result
}

// RenderSettings renders the settings menu
func RenderSettings(cm *ConfigManager, width, height int, wantToQuit bool) string {
	var s strings.Builder

	title := titleStyle.Render("settings")
	s.WriteString(title)
	s.WriteString("\n\n")

	s.WriteString(subtleStyle.Render("configuration options:"))
	s.WriteString("\n\n")

	s.WriteString("   " + wpmStyle.Render("1") + subtleStyle.Render(" → cursor style: ") +
		accuracyStyle.Render(cm.GetCursorTypeName()))
	s.WriteString("\n")
	s.WriteString("   " + wpmStyle.Render("2") + subtleStyle.Render(" → accent color: ") +
		lipgloss.NewStyle().Foreground(lipgloss.Color(cm.GetAccentColorHex())).
			Bold(true).Render(cm.GetAccentColorName()))
	s.WriteString("\n\n")

	s.WriteString(subtleStyle.Render("current settings:"))
	s.WriteString("\n")
	s.WriteString(fmt.Sprintf("   cursor: %s\n", cm.GetCursorTypeName()))
	s.WriteString(fmt.Sprintf("   color:  %s (%s)\n", cm.GetAccentColorName(), cm.GetAccentColorHex()))

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

// RenderCursorSelect renders the cursor type selection screen
func RenderCursorSelect(cm *ConfigManager, width, height int, wantToQuit bool) string {
	var s strings.Builder

	title := titleStyle.Render("cursor style")
	s.WriteString(title)
	s.WriteString("\n\n")

	s.WriteString(subtleStyle.Render("select cursor type:"))
	s.WriteString("\n\n")

	currentCursor := cm.GetCursorType()
	cursorTypes := []struct {
		key   string
		name  string
		ctype CursorType
	}{
		{"1", "Block", CursorBlock},
		{"2", "Line", CursorLine},
		{"3", "Underline", CursorUnderline},
		{"4", "Bar", CursorBar},
	}

	for _, ct := range cursorTypes {
		keyStyle := wpmStyle
		labelStyle := subtleStyle
		if ct.ctype == currentCursor {
			keyStyle = lipgloss.NewStyle().Foreground(colorAccent).Bold(true)
			labelStyle = lipgloss.NewStyle().Foreground(colorText)
		}
		s.WriteString(fmt.Sprintf("   %s %s\n", keyStyle.Render(ct.key), labelStyle.Render("→ "+ct.name)))
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

// RenderColorSelect renders the accent color selection screen
func RenderColorSelect(cm *ConfigManager, width, height int, wantToQuit bool) string {
	var s strings.Builder

	title := titleStyle.Render("accent color")
	s.WriteString(title)
	s.WriteString("\n\n")

	s.WriteString(subtleStyle.Render("select a color:"))
	s.WriteString("\n\n")

	colors := []struct {
		key   string
		name  string
		color string
	}{
		{"1", "Red", "#e06c75"},
		{"2", "Orange", "#d19a66"},
		{"3", "Yellow", "#e5c07b"},
		{"4", "Green", "#98c379"},
		{"5", "Cyan", "#56b6c2"},
		{"6", "Blue", "#61afef"},
		{"7", "Purple", "#c678dd"},
		{"8", "Pink", "#ff79c6"},
		{"9", "White", "#abb2bf"},
		{"0", "Black", "#282c34"},
	}

	currentHex := cm.GetAccentColorHex()

	for _, c := range colors {
		keyStyle := wpmStyle
		nameStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(c.color))
		if c.color == currentHex {
			keyStyle = lipgloss.NewStyle().Foreground(colorAccent).Bold(true)
			nameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(c.color)).Bold(true)
		}
		s.WriteString(fmt.Sprintf("   %s %s\n", keyStyle.Render(c.key), nameStyle.Render("● "+c.name)))
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

// RenderChallenges renders the daily challenges screen
func RenderChallenges(dc *DailyChallenges, width, height int, wantToQuit bool) string {
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
