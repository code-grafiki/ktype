package ui

import (
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
