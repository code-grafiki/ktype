package app

import (
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"ktype/internal/game"
	"ktype/internal/storage"
	"ktype/internal/ui"
	"ktype/internal/words"
)

// Update handles all messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case tickMsg:
		// Clear quit prompt after 2 seconds
		if m.WantToQuit && time.Since(m.QuitPressAt) > 2*time.Second {
			m.WantToQuit = false
		}

		if m.Game != nil && m.State == game.StatePlaying {
			m.Game.Update()
			if m.Game.State == game.StateFinished {
				m.State = game.StateFinished
				// Save score to leaderboard
				m.Leaderboard.AddScore(m.Game.WPM(), m.Game.Accuracy(), m.Game.ModeString())
				// Update challenges progress
				m.Challenges.UpdateProgress(m.Game.WPM(), m.Game.Accuracy(), len(m.Game.TypedWords))
				return m, nil
			}
			return m, tickCmd()
		}
		return m, nil
	}

	return m, nil
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle ctrl+c globally
	if msg.String() == "ctrl+c" {
		return m, tea.Quit
	}

	// Reset WantToQuit on any other key
	if msg.Type != tea.KeyEsc {
		m.WantToQuit = false
	}

	switch m.State {
	case game.StateMenu:
		return m.handleMenuKey(msg)
	case game.StateDifficultySelect:
		return m.handleDifficultySelectKey(msg)
	case game.StateComplexitySelect:
		return m.handleComplexitySelectKey(msg)
	case game.StateStats:
		return m.handleStatsKey(msg)
	case game.StateHeatmap:
		return m.handleHeatmapKey(msg)
	case game.StateSettings:
		return m.handleSettingsKey(msg)
	case game.StateCursorSelect:
		return m.handleCursorSelectKey(msg)
	case game.StateColorSelect:
		return m.handleColorSelectKey(msg)
	case game.StateCustomWordList:
		return m.handleCustomWordListKey(msg)
	case game.StateTimeSelect:
		return m.handleTimeSelectKey(msg)
	case game.StateWordsSelect:
		return m.handleWordsSelectKey(msg)
	case game.StateCustomInput:
		return m.handleCustomInputKey(msg)
	case game.StatePlaying:
		return m.handlePlayingKey(msg)
	case game.StateFinished:
		return m.handleFinishedKey(msg)
	case game.StateChallenges:
		return m.handleChallengesKey(msg)
	}
	return m, nil
}

func (m Model) handleMenuKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "1": // Quick start 30s
		m.Game = game.NewTimed(30*time.Second, m.Difficulty, m.Complexity, m.Heatmap)
		m.State = game.StatePlaying
		return m, tickCmd()
	case "2": // Quick start 50 words
		m.Game = game.NewWords(50, m.Difficulty, m.Complexity, m.Heatmap)
		m.State = game.StatePlaying
		return m, tickCmd()
	case "3": // Zen mode
		m.Game = game.NewZen(m.Difficulty, m.Complexity, m.Heatmap)
		m.State = game.StatePlaying
		return m, tickCmd()
	case "t":
		m.State = game.StateTimeSelect
		return m, nil
	case "w":
		m.State = game.StateWordsSelect
		return m, nil
	case "d":
		m.State = game.StateDifficultySelect
		return m, nil
	case "c":
		m.State = game.StateComplexitySelect
		return m, nil
	case "s":
		m.State = game.StateStats
		return m, nil
	case "h":
		m.State = game.StateHeatmap
		return m, nil
	case "l":
		m.State = game.StateCustomWordList
		return m, nil
	case "v":
		m.State = game.StateChallenges
		return m, nil
	case ",":
		m.State = game.StateSettings
		return m, nil
	case "esc":
		if m.WantToQuit {
			return m, tea.Quit
		}
		m.WantToQuit = true
		m.QuitPressAt = time.Now()
		return m, nil
	}
	return m, nil
}

func (m Model) handleDifficultySelectKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "1":
		m.Difficulty = words.DifficultyEasy
		m.State = game.StateMenu
		return m, nil
	case "2":
		m.Difficulty = words.DifficultyMedium
		m.State = game.StateMenu
		return m, nil
	case "3":
		m.Difficulty = words.DifficultyHard
		m.State = game.StateMenu
		return m, nil
	case "esc":
		m.State = game.StateMenu
		return m, nil
	}
	return m, nil
}

func (m Model) handleComplexitySelectKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "1":
		m.Complexity = words.ComplexityNormal
		m.State = game.StateMenu
		return m, nil
	case "2":
		m.Complexity = words.ComplexityPunctuation
		m.State = game.StateMenu
		return m, nil
	case "3":
		m.Complexity = words.ComplexityNumbers
		m.State = game.StateMenu
		return m, nil
	case "4":
		m.Complexity = words.ComplexityFull
		m.State = game.StateMenu
		return m, nil
	case "esc":
		m.State = game.StateMenu
		return m, nil
	}
	return m, nil
}

func (m Model) handleStatsKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.State = game.StateMenu
		return m, nil
	}
	return m, nil
}

func (m Model) handleHeatmapKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.State = game.StateMenu
		return m, nil
	case "r":
		if m.Heatmap != nil {
			m.Heatmap.Clear()
		}
		return m, nil
	}
	return m, nil
}

func (m Model) handleCustomWordListKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.State = game.StateMenu
		return m, nil
	case "1", "2", "3", "4", "5", "6", "7", "8", "9":
		// Set current word list based on selection
		index, _ := strconv.Atoi(msg.String())
		lists := m.WordListManager.ListNames()
		if index > 0 && index <= len(lists) {
			m.CurrentWordList = lists[index-1]
		}
		return m, nil
	case "d":
		// Delete the selected word list
		if m.CurrentWordList != "" {
			m.WordListManager.DeleteList(m.CurrentWordList)
			m.CurrentWordList = ""
		}
		return m, nil
	}
	return m, nil
}

func (m Model) handleTimeSelectKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "1":
		m.Game = game.NewTimed(15*time.Second, m.Difficulty, m.Complexity, m.Heatmap)
		m.State = game.StatePlaying
		return m, tickCmd()
	case "2":
		m.Game = game.NewTimed(30*time.Second, m.Difficulty, m.Complexity, m.Heatmap)
		m.State = game.StatePlaying
		return m, tickCmd()
	case "3":
		m.Game = game.NewTimed(60*time.Second, m.Difficulty, m.Complexity, m.Heatmap)
		m.State = game.StatePlaying
		return m, tickCmd()
	case "c":
		m.CustomInput = ""
		m.InputMode = "time"
		m.State = game.StateCustomInput
		m.WantToQuit = false
		return m, nil
	case "esc":
		m.State = game.StateMenu
		return m, nil
	}
	return m, nil
}

func (m Model) handleWordsSelectKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "1":
		m.Game = game.NewWords(10, m.Difficulty, m.Complexity, m.Heatmap)
		m.State = game.StatePlaying
		return m, tickCmd()
	case "2":
		m.Game = game.NewWords(25, m.Difficulty, m.Complexity, m.Heatmap)
		m.State = game.StatePlaying
		return m, tickCmd()
	case "3":
		m.Game = game.NewWords(50, m.Difficulty, m.Complexity, m.Heatmap)
		m.State = game.StatePlaying
		return m, tickCmd()
	case "4":
		m.Game = game.NewWords(100, m.Difficulty, m.Complexity, m.Heatmap)
		m.State = game.StatePlaying
		return m, tickCmd()
	case "c":
		m.CustomInput = ""
		m.InputMode = "words"
		m.State = game.StateCustomInput
		m.WantToQuit = false
		return m, nil
	case "esc":
		m.State = game.StateMenu
		return m, nil
	}
	return m, nil
}

func (m Model) handleCustomInputKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.CustomInput = ""
		if m.InputMode == "time" {
			m.State = game.StateTimeSelect
		} else {
			m.State = game.StateWordsSelect
		}
		return m, nil
	case tea.KeyEnter:
		if len(m.CustomInput) > 0 {
			value, err := strconv.Atoi(m.CustomInput)
			if err == nil && value > 0 {
				// Validate bounds
				if m.InputMode == "time" && value > 3600 {
					return m, nil // Max 1 hour
				}
				if m.InputMode == "words" && value > 1000 {
					return m, nil // Max 1000 words
				}
				if m.InputMode == "time" {
					m.Game = game.NewTimed(time.Duration(value)*time.Second, m.Difficulty, m.Complexity, m.Heatmap)
				} else {
					m.Game = game.NewWords(value, m.Difficulty, m.Complexity, m.Heatmap)
				}
				m.State = game.StatePlaying
				m.CustomInput = ""
				return m, tickCmd()
			}
		}
		return m, nil
	case tea.KeyBackspace:
		if len(m.CustomInput) > 0 {
			m.CustomInput = m.CustomInput[:len(m.CustomInput)-1]
		}
		return m, nil
	case tea.KeyRunes:
		// Only allow digits
		for _, r := range msg.Runes {
			if r >= '0' && r <= '9' && len(m.CustomInput) < 4 {
				m.CustomInput += string(r)
			}
		}
		return m, nil
	}
	return m, nil
}

func (m Model) handlePlayingKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		if m.WantToQuit {
			m.Game = nil
			m.State = game.StateMenu
			m.WantToQuit = false
			return m, nil
		}
		m.WantToQuit = true
		m.QuitPressAt = time.Now()
		return m, nil
	case "tab":
		// Restart - go back to menu
		m.Game = nil
		m.State = game.StateMenu
		m.WantToQuit = false
		return m, nil
	}

	switch msg.Type {
	case tea.KeyBackspace:
		m.WantToQuit = false
		if m.Game != nil {
			m.Game.HandleBackspace()
		}
		return m, nil
	case tea.KeySpace:
		m.WantToQuit = false
		if m.Game != nil {
			m.Game.HandleSpace()
			// Check if game finished (words mode)
			if m.Game.State == game.StateFinished {
				m.State = game.StateFinished
				m.Leaderboard.AddScore(m.Game.WPM(), m.Game.Accuracy(), m.Game.ModeString())
			}
		}
		return m, nil
	case tea.KeyRunes:
		m.WantToQuit = false
		if m.Game != nil && len(msg.Runes) > 0 {
			for _, r := range msg.Runes {
				m.Game.HandleChar(r)
			}
		}
		return m, nil
	}
	return m, nil
}

func (m Model) handleFinishedKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "tab", "enter":
		m.Game = nil
		m.State = game.StateMenu
		m.WantToQuit = false
		return m, nil
	case "esc":
		if m.WantToQuit {
			return m, tea.Quit
		}
		m.WantToQuit = true
		m.QuitPressAt = time.Now()
		return m, nil
	}
	return m, nil
}

func (m Model) handleChallengesKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.State = game.StateMenu
		return m, nil
	}
	return m, nil
}

func (m Model) handleSettingsKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.State = game.StateMenu
		return m, nil
	case "1":
		m.State = game.StateCursorSelect
		return m, nil
	case "2":
		m.State = game.StateColorSelect
		return m, nil
	}
	return m, nil
}

func (m Model) handleCursorSelectKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.State = game.StateSettings
		return m, nil
	case "1":
		m.ConfigManager.SetCursorType(storage.CursorBlock)
		m.State = game.StateSettings
		return m, nil
	case "2":
		m.ConfigManager.SetCursorType(storage.CursorLine)
		m.State = game.StateSettings
		return m, nil
	case "3":
		m.ConfigManager.SetCursorType(storage.CursorUnderline)
		m.State = game.StateSettings
		return m, nil
	case "4":
		m.ConfigManager.SetCursorType(storage.CursorBar)
		m.State = game.StateSettings
		return m, nil
	}
	return m, nil
}

func (m Model) handleColorSelectKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "1":
		m.ConfigManager.SetAccentColor(storage.ColorRed)
		ui.UpdateAccentColor(m.ConfigManager.GetAccentColor())
		m.State = game.StateSettings
		return m, nil
	case "2":
		m.ConfigManager.SetAccentColor(storage.ColorOrange)
		ui.UpdateAccentColor(m.ConfigManager.GetAccentColor())
		m.State = game.StateSettings
		return m, nil
	case "3":
		m.ConfigManager.SetAccentColor(storage.ColorYellow)
		ui.UpdateAccentColor(m.ConfigManager.GetAccentColor())
		m.State = game.StateSettings
		return m, nil
	case "4":
		m.ConfigManager.SetAccentColor(storage.ColorGreen)
		ui.UpdateAccentColor(m.ConfigManager.GetAccentColor())
		m.State = game.StateSettings
		return m, nil
	case "5":
		m.ConfigManager.SetAccentColor(storage.ColorCyan)
		ui.UpdateAccentColor(m.ConfigManager.GetAccentColor())
		m.State = game.StateSettings
		return m, nil
	case "6":
		m.ConfigManager.SetAccentColor(storage.ColorBlue)
		ui.UpdateAccentColor(m.ConfigManager.GetAccentColor())
		m.State = game.StateSettings
		return m, nil
	case "7":
		m.ConfigManager.SetAccentColor(storage.ColorPurple)
		ui.UpdateAccentColor(m.ConfigManager.GetAccentColor())
		m.State = game.StateSettings
		return m, nil
	case "8":
		m.ConfigManager.SetAccentColor(storage.ColorPink)
		ui.UpdateAccentColor(m.ConfigManager.GetAccentColor())
		m.State = game.StateSettings
		return m, nil
	case "9":
		m.ConfigManager.SetAccentColor(storage.ColorWhite)
		ui.UpdateAccentColor(m.ConfigManager.GetAccentColor())
		m.State = game.StateSettings
		return m, nil
	case "0":
		m.ConfigManager.SetAccentColor(storage.ColorBlack)
		ui.UpdateAccentColor(m.ConfigManager.GetAccentColor())
		m.State = game.StateSettings
		return m, nil
	case "esc":
		m.State = game.StateSettings
		return m, nil
	}
	return m, nil
}
