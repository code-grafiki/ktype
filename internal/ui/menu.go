package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"ktype/internal/storage"
	"ktype/internal/words"
)

// RenderMainMenu renders the main menu with quick start options
func RenderMainMenu(lb *storage.Leaderboard, width, height int, wantToQuit bool, difficulty words.Difficulty, complexity words.Complexity) string {
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

	formatPB := func(s *storage.Score) string {
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
func RenderTimeSelect(lb *storage.Leaderboard, width, height int, wantToQuit bool) string {
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
func RenderWordsSelect(lb *storage.Leaderboard, width, height int, wantToQuit bool) string {
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
func RenderDifficultySelect(currentDifficulty words.Difficulty, width, height int, wantToQuit bool) string {
	var s strings.Builder

	title := titleStyle.Render("difficulty")
	s.WriteString(title)
	s.WriteString("\n\n")

	s.WriteString(subtleStyle.Render("select difficulty:"))
	s.WriteString("\n\n")

	options := []struct {
		key        string
		label      string
		difficulty words.Difficulty
		desc       string
	}{
		{"1", "easy", words.DifficultyEasy, "2-4 letter words"},
		{"2", "medium", words.DifficultyMedium, "5-7 letter words"},
		{"3", "hard", words.DifficultyHard, "8+ letter words"},
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
func RenderComplexitySelect(currentComplexity words.Complexity, width, height int, wantToQuit bool) string {
	var s strings.Builder

	title := titleStyle.Render("complexity")
	s.WriteString(title)
	s.WriteString("\n\n")

	s.WriteString(subtleStyle.Render("select complexity:"))
	s.WriteString("\n\n")

	options := []struct {
		key        string
		label      string
		complexity words.Complexity
		desc       string
	}{
		{"1", "normal", words.ComplexityNormal, "letters only"},
		{"2", "punctuation", words.ComplexityPunctuation, "letters + punctuation"},
		{"3", "numbers", words.ComplexityNumbers, "letters + numbers"},
		{"4", "full", words.ComplexityFull, "letters + punctuation + numbers"},
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

// RenderSettings renders the settings menu
func RenderSettings(cm *storage.ConfigManager, width, height int, wantToQuit bool) string {
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
func RenderCursorSelect(cm *storage.ConfigManager, width, height int, wantToQuit bool) string {
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
		ctype storage.CursorType
	}{
		{"1", "Block", storage.CursorBlock},
		{"2", "Line", storage.CursorLine},
		{"3", "Underline", storage.CursorUnderline},
		{"4", "Bar", storage.CursorBar},
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
func RenderColorSelect(cm *storage.ConfigManager, width, height int, wantToQuit bool) string {
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

// RenderCustomWordList renders the custom word list management screen
func RenderCustomWordList(wm *storage.WordListManager, currentList string, width, height int, wantToQuit bool) string {
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
