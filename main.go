package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// tickMsg is sent every tick to update the timer
type tickMsg time.Time

// model is the Bubble Tea model
type model struct {
	game        *Game
	leaderboard *Leaderboard
	width       int
	height      int
	state       GameState
	wantToQuit  bool      // For double-esc to quit
	quitPressAt time.Time // When first esc was pressed

	// For custom input
	customInput string
	inputMode   string         // "time" or "words"
	difficulty  Difficulty     // Current difficulty level
	complexity  WordComplexity // Current complexity level (normal, punctuation, numbers, full)

	// For custom word lists
	wordListManager *WordListManager
	currentWordList string // name of selected word list

	// For heatmap (persists across games)
	heatmap *Heatmap

	// For configuration
	configManager *ConfigManager

	// Daily challenges
	challenges *DailyChallenges
}

func initialModel() model {
	return model{
		state:           StateMenu,
		width:           80,
		height:          24,
		leaderboard:     NewLeaderboard(),
		difficulty:      DifficultyMedium,
		complexity:      ComplexityNormal,
		wordListManager: NewWordListManager(),
		currentWordList: "",
		heatmap:         NewHeatmap(),
		configManager:   NewConfigManager(),
		challenges:      NewDailyChallenges(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tickMsg:
		// Clear quit prompt after 2 seconds
		if m.wantToQuit && time.Since(m.quitPressAt) > 2*time.Second {
			m.wantToQuit = false
		}

		if m.game != nil && m.state == StatePlaying {
			m.game.Update()
			if m.game.State == StateFinished {
				m.state = StateFinished
				// Save score to leaderboard
				m.leaderboard.AddScore(m.game.WPM(), m.game.Accuracy(), m.game.ModeString())
				// Update challenges progress
				m.challenges.UpdateProgress(m.game.WPM(), m.game.Accuracy(), len(m.game.TypedWords))
				return m, nil
			}
			return m, tickCmd()
		}
		return m, nil
	}

	return m, nil
}

func (m model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle ctrl+c globally
	if msg.String() == "ctrl+c" {
		return m, tea.Quit
	}

	// Reset wantToQuit on any other key
	if msg.Type != tea.KeyEsc {
		m.wantToQuit = false
	}

	switch m.state {
	case StateMenu:
		return m.handleMenuKey(msg)
	case StateDifficultySelect:
		return m.handleDifficultySelectKey(msg)
	case StateComplexitySelect:
		return m.handleComplexitySelectKey(msg)
	case StateStats:
		return m.handleStatsKey(msg)
	case StateHeatmap:
		return m.handleHeatmapKey(msg)
	case StateSettings:
		return m.handleSettingsKey(msg)
	case StateCursorSelect:
		return m.handleCursorSelectKey(msg)
	case StateColorSelect:
		return m.handleColorSelectKey(msg)
	case StateCustomWordList:
		return m.handleCustomWordListKey(msg)
	case StateTimeSelect:
		return m.handleTimeSelectKey(msg)
	case StateWordsSelect:
		return m.handleWordsSelectKey(msg)
	case StateCustomInput:
		return m.handleCustomInputKey(msg)
	case StatePlaying:
		return m.handlePlayingKey(msg)
	case StateFinished:
		return m.handleFinishedKey(msg)
	case StateChallenges:
		return m.handleChallengesKey(msg)
	}
	return m, nil
}

func (m model) handleMenuKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "1": // Quick start 30s
		m.game = NewTimedGame(30*time.Second, m.difficulty, m.complexity, m.heatmap)
		m.state = StatePlaying
		return m, tickCmd()
	case "2": // Quick start 50 words
		m.game = NewWordsGame(50, m.difficulty, m.complexity, m.heatmap)
		m.state = StatePlaying
		return m, tickCmd()
	case "3": // Zen mode
		m.game = NewZenGame(m.difficulty, m.complexity, m.heatmap)
		m.state = StatePlaying
		return m, tickCmd()
	case "t":
		m.state = StateTimeSelect
		return m, nil
	case "w":
		m.state = StateWordsSelect
		return m, nil
	case "d":
		m.state = StateDifficultySelect
		return m, nil
	case "c":
		m.state = StateComplexitySelect
		return m, nil
	case "s":
		m.state = StateStats
		return m, nil
	case "h":
		m.state = StateHeatmap
		return m, nil
	case "l":
		m.state = StateCustomWordList
		return m, nil
	case "v":
		m.state = StateChallenges
		return m, nil
	case ",":
		m.state = StateSettings
		return m, nil
	case "esc":
		if m.wantToQuit {
			return m, tea.Quit
		}
		m.wantToQuit = true
		m.quitPressAt = time.Now()
		return m, nil
	}
	return m, nil
}

func (m model) handleDifficultySelectKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "1":
		m.difficulty = DifficultyEasy
		m.state = StateMenu
		return m, nil
	case "2":
		m.difficulty = DifficultyMedium
		m.state = StateMenu
		return m, nil
	case "3":
		m.difficulty = DifficultyHard
		m.state = StateMenu
		return m, nil
	case "esc":
		m.state = StateMenu
		return m, nil
	}
	return m, nil
}

func (m model) handleComplexitySelectKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "1":
		m.complexity = ComplexityNormal
		m.state = StateMenu
		return m, nil
	case "2":
		m.complexity = ComplexityPunctuation
		m.state = StateMenu
		return m, nil
	case "3":
		m.complexity = ComplexityNumbers
		m.state = StateMenu
		return m, nil
	case "4":
		m.complexity = ComplexityFull
		m.state = StateMenu
		return m, nil
	case "esc":
		m.state = StateMenu
		return m, nil
	}
	return m, nil
}

func (m model) handleStatsKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.state = StateMenu
		return m, nil
	}
	return m, nil
}

func (m model) handleHeatmapKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.state = StateMenu
		return m, nil
	case "r":
		if m.heatmap != nil {
			m.heatmap.Clear()
		}
		return m, nil
	}
	return m, nil
}

func (m model) handleCustomWordListKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.state = StateMenu
		return m, nil
	case "1", "2", "3", "4", "5", "6", "7", "8", "9":
		// Set current word list based on selection
		index, _ := strconv.Atoi(msg.String())
		lists := m.wordListManager.ListNames()
		if index > 0 && index <= len(lists) {
			m.currentWordList = lists[index-1]
		}
		return m, nil
	case "d":
		// Delete the selected word list
		if m.currentWordList != "" {
			m.wordListManager.DeleteList(m.currentWordList)
			m.currentWordList = ""
		}
		return m, nil
	case "n":
		// Prompt for new word list creation
		fmt.Println("New word list creation - press enter to continue")
		return m, nil
	}
	return m, nil
}

func (m model) handleTimeSelectKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "1":
		m.game = NewTimedGame(15*time.Second, m.difficulty, m.complexity, m.heatmap)
		m.state = StatePlaying
		return m, tickCmd()
	case "2":
		m.game = NewTimedGame(30*time.Second, m.difficulty, m.complexity, m.heatmap)
		m.state = StatePlaying
		return m, tickCmd()
	case "3":
		m.game = NewTimedGame(60*time.Second, m.difficulty, m.complexity, m.heatmap)
		m.state = StatePlaying
		return m, tickCmd()
	case "c":
		m.customInput = ""
		m.inputMode = "time"
		m.state = StateCustomInput
		m.wantToQuit = false
		return m, nil
	case "esc":
		m.state = StateMenu
		return m, nil
	}
	return m, nil
}

func (m model) handleWordsSelectKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "1":
		m.game = NewWordsGame(10, m.difficulty, m.complexity, m.heatmap)
		m.state = StatePlaying
		return m, tickCmd()
	case "2":
		m.game = NewWordsGame(25, m.difficulty, m.complexity, m.heatmap)
		m.state = StatePlaying
		return m, tickCmd()
	case "3":
		m.game = NewWordsGame(50, m.difficulty, m.complexity, m.heatmap)
		m.state = StatePlaying
		return m, tickCmd()
	case "4":
		m.game = NewWordsGame(100, m.difficulty, m.complexity, m.heatmap)
		m.state = StatePlaying
		return m, tickCmd()
	case "c":
		m.customInput = ""
		m.inputMode = "words"
		m.state = StateCustomInput
		m.wantToQuit = false
		return m, nil
	case "esc":
		m.state = StateMenu
		return m, nil
	}
	return m, nil
}

func (m model) handleCustomInputKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.customInput = ""
		if m.inputMode == "time" {
			m.state = StateTimeSelect
		} else {
			m.state = StateWordsSelect
		}
		return m, nil
	case tea.KeyEnter:
		if len(m.customInput) > 0 {
			value, err := strconv.Atoi(m.customInput)
			if err == nil && value > 0 {
				// Validate bounds
				if m.inputMode == "time" && value > 3600 {
					return m, nil // Max 1 hour
				}
				if m.inputMode == "words" && value > 1000 {
					return m, nil // Max 1000 words
				}
				if m.inputMode == "time" {
					m.game = NewTimedGame(time.Duration(value)*time.Second, m.difficulty, m.complexity, m.heatmap)
				} else {
					m.game = NewWordsGame(value, m.difficulty, m.complexity, m.heatmap)
				}
				m.state = StatePlaying
				m.customInput = ""
				return m, tickCmd()
			}
		}
		return m, nil
	case tea.KeyBackspace:
		if len(m.customInput) > 0 {
			m.customInput = m.customInput[:len(m.customInput)-1]
		}
		return m, nil
	case tea.KeyRunes:
		// Only allow digits
		for _, r := range msg.Runes {
			if r >= '0' && r <= '9' && len(m.customInput) < 4 {
				m.customInput += string(r)
			}
		}
		return m, nil
	}
	return m, nil
}

func (m model) handlePlayingKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		if m.wantToQuit {
			m.game = nil
			m.state = StateMenu
			m.wantToQuit = false
			return m, nil
		}
		m.wantToQuit = true
		m.quitPressAt = time.Now()
		return m, nil
	case "tab":
		// Restart - go back to menu
		m.game = nil
		m.state = StateMenu
		m.wantToQuit = false
		return m, nil
	}

	switch msg.Type {
	case tea.KeyBackspace:
		m.wantToQuit = false
		if m.game != nil {
			m.game.HandleBackspace()
		}
		return m, nil
	case tea.KeySpace:
		m.wantToQuit = false
		if m.game != nil {
			m.game.HandleSpace()
			// Check if game finished (words mode)
			if m.game.State == StateFinished {
				m.state = StateFinished
				m.leaderboard.AddScore(m.game.WPM(), m.game.Accuracy(), m.game.ModeString())
			}
		}
		return m, nil
	case tea.KeyRunes:
		m.wantToQuit = false
		if m.game != nil && len(msg.Runes) > 0 {
			for _, r := range msg.Runes {
				m.game.HandleChar(r)
			}
		}
		return m, nil
	}
	return m, nil
}

func (m model) handleFinishedKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "tab", "enter":
		m.game = nil
		m.state = StateMenu
		m.wantToQuit = false
		return m, nil
	case "esc":
		if m.wantToQuit {
			return m, tea.Quit
		}
		m.wantToQuit = true
		m.quitPressAt = time.Now()
		return m, nil
	}
	return m, nil
}

func (m model) handleChallengesKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.state = StateMenu
		return m, nil
	}
	return m, nil
}

func (m model) handleSettingsKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.state = StateMenu
		return m, nil
	case "1":
		m.state = StateCursorSelect
		return m, nil
	case "2":
		m.state = StateColorSelect
		return m, nil
	}
	return m, nil
}

func (m model) handleCursorSelectKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.state = StateSettings
		return m, nil
	case "1":
		m.configManager.SetCursorType(CursorBlock)
		m.state = StateSettings
		return m, nil
	case "2":
		m.configManager.SetCursorType(CursorLine)
		m.state = StateSettings
		return m, nil
	case "3":
		m.configManager.SetCursorType(CursorUnderline)
		m.state = StateSettings
		return m, nil
	case "4":
		m.configManager.SetCursorType(CursorBar)
		m.state = StateSettings
		return m, nil
	}
	return m, nil
}

func (m model) handleColorSelectKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.customInput = ""
		m.state = StateSettings
		return m, nil
	case tea.KeyEnter:
		if len(m.customInput) > 0 {
			// Check if it's a preset number
			switch m.customInput {
			case "1":
				m.configManager.SetAccentColor(ColorRed)
			case "2":
				m.configManager.SetAccentColor(ColorOrange)
			case "3":
				m.configManager.SetAccentColor(ColorYellow)
			case "4":
				m.configManager.SetAccentColor(ColorGreen)
			case "5":
				m.configManager.SetAccentColor(ColorCyan)
			case "6":
				m.configManager.SetAccentColor(ColorBlue)
			case "7":
				m.configManager.SetAccentColor(ColorPurple)
			case "8":
				m.configManager.SetAccentColor(ColorPink)
			case "9":
				m.configManager.SetAccentColor(ColorWhite)
			case "0":
				m.configManager.SetAccentColor(ColorBlack)
			default:
				// Try to parse as hex color
				if m.configManager.SetCustomColor(m.customInput) {
					m.configManager.SetAccentColor(ColorCustom)
				}
			}
			m.customInput = ""
			m.state = StateSettings
		}
		return m, nil
	case tea.KeyBackspace:
		if len(m.customInput) > 0 {
			m.customInput = m.customInput[:len(m.customInput)-1]
		}
		return m, nil
	case tea.KeyRunes:
		// Allow hex characters (0-9, a-f, A-F) and #
		for _, r := range msg.Runes {
			if (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F') || r == '#' {
				if len(m.customInput) < 7 {
					m.customInput += string(r)
				}
			}
		}
		return m, nil
	}
	return m, nil
}

func (m model) View() string {
	switch m.state {
	case StateMenu:
		return RenderMainMenu(m.leaderboard, m.width, m.height, m.wantToQuit, m.difficulty, m.complexity)
	case StateDifficultySelect:
		return RenderDifficultySelect(m.difficulty, m.width, m.height, m.wantToQuit)
	case StateComplexitySelect:
		return RenderComplexitySelect(m.complexity, m.width, m.height, m.wantToQuit)
	case StateStats:
		return RenderStats(NewStatistics(m.leaderboard), m.width, m.height, m.wantToQuit)
	case StateHeatmap:
		return RenderHeatmap(m.heatmap, m.width, m.height, m.wantToQuit)
	case StateSettings:
		return RenderSettings(m.configManager, m.width, m.height, m.wantToQuit)
	case StateCursorSelect:
		return RenderCursorSelect(m.configManager, m.width, m.height, m.wantToQuit)
	case StateColorSelect:
		return RenderColorSelect(m.configManager, m.width, m.height, m.wantToQuit)
	case StateCustomWordList:
		return RenderCustomWordList(m.wordListManager, m.currentWordList, m.width, m.height, m.wantToQuit)
	case StateTimeSelect:
		return RenderTimeSelect(m.leaderboard, m.width, m.height, m.wantToQuit)
	case StateWordsSelect:
		return RenderWordsSelect(m.leaderboard, m.width, m.height, m.wantToQuit)
	case StateCustomInput:
		return RenderCustomInput(m.customInput, m.inputMode, m.width, m.height, m.configManager.GetCursorType().CursorChar())
	case StatePlaying:
		if m.game != nil {
			return RenderGame(m.game, m.width, m.height, m.wantToQuit)
		}
	case StateFinished:
		if m.game != nil {
			isPB := m.leaderboard.IsPB(m.game.WPM(), m.game.ModeString())
			return RenderFinished(m.game, m.width, m.height, isPB, m.wantToQuit)
		}
	case StateChallenges:
		return RenderChallenges(m.challenges, m.width, m.height, m.wantToQuit)
	}
	return ""
}

func tickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func main() {
	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
