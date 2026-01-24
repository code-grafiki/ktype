package main

import (
	"fmt"
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
			Bold(true)

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

// RenderMainMenu renders the main menu with quick start options
func RenderMainMenu(lb *Leaderboard, width, height int, wantToQuit bool) string {
	var s strings.Builder

	title := titleStyle.Render(" ⌨  ktype")
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

	title := titleStyle.Render(" ⌨  timed mode")
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

	title := titleStyle.Render(" ⌨  words mode")
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

// RenderCustomInput renders the custom input screen
func RenderCustomInput(input, mode string, width, height int) string {
	var s strings.Builder

	var prompt string
	if mode == "time" {
		prompt = "enter duration (seconds)"
	} else {
		prompt = "enter word count"
	}

	title := titleStyle.Render(" ⌨  " + prompt)
	s.WriteString(title)
	s.WriteString("\n\n")

	// Show input with cursor
	inputDisplay := wpmStyle.Render(input) + cursorStyle.Render("▌")
	s.WriteString("        " + inputDisplay)
	s.WriteString("\n\n")

	s.WriteString(helpStyle.Render("enter to confirm • esc to go back"))

	content := containerStyle.Render(s.String())
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}

// RenderGame renders the main game screen
func RenderGame(g *Game, width, height int, wantToQuit bool) string {
	var s strings.Builder

	// Inside the container, we have Width(65) and Padding(2, 4).
	// So internal width is 65 - 4*2 = 57.
	internalWidth := 57

	// Build the words display (3 lines)
	// We use a slightly smaller width for the words themselves to ensure they fit well
	wordsLines := buildWordsLines(g, internalWidth-2, 3)

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
	stats := statsStyle.Render(wpm + subtleStyle.Render("  •  ") + accuracy)
	s.WriteString(lipgloss.PlaceHorizontal(internalWidth, lipgloss.Center, stats))

	// Help text or quit confirmation
	s.WriteString("\n\n")
	var help string
	if wantToQuit {
		help = errorStyle.Render("press esc again to quit")
	} else {
		help = helpStyle.Render("r/tab: restart • esc: abort")
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
		title := titleStyle.Render(" ⌨  test complete!")
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
		help = helpStyle.Render("r/tab to restart • esc to quit")
	}
	s.WriteString(help)

	content := containerStyle.Render(s.String())
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}

// buildWordsLines builds 3 lines of words with proper scrolling
func buildWordsLines(g *Game, maxWidth int, numLines int) []string {
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
				cursor := cursorStyle.Render("▌")
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
