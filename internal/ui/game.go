package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"ktype/internal/game"
)

// RenderGame renders the main game screen
func RenderGame(g *game.Game, width, height int, wantToQuit bool, cursorChar string) string {
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
func RenderFinished(g *game.Game, width, height int, isPB bool, wantToQuit bool) string {
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
func buildWordsLines(g *game.Game, maxWidth int, numLines int, cursorChar string) []string {
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

// GetErrorStats returns statistics about errors in the current game
func GetErrorStats(g *game.Game) (totalErrors int, topErrors []struct {
	Char  rune
	Count int
}) {
	totalErrors = len(g.Errors)

	// Count errors by character
	charCounts := make(map[rune]int)
	for _, err := range g.Errors {
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
