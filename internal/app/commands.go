package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// tickCmd returns a command that ticks every 100ms
func tickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
