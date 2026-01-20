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
	inputMode   string // "time" or "words"
}

func initialModel() model {
	return model{
		state:       StateMenu,
		width:       80,
		height:      24,
		leaderboard: NewLeaderboard(),
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
	}
	return m, nil
}

func (m model) handleMenuKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "1": // Quick start 30s
		m.game = NewTimedGame(30 * time.Second)
		m.state = StatePlaying
		return m, tickCmd()
	case "2": // Quick start 50 words
		m.game = NewWordsGame(50)
		m.state = StatePlaying
		return m, tickCmd()
	case "t":
		m.state = StateTimeSelect
		return m, nil
	case "w":
		m.state = StateWordsSelect
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

func (m model) handleTimeSelectKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "1":
		m.game = NewTimedGame(15 * time.Second)
		m.state = StatePlaying
		return m, tickCmd()
	case "2":
		m.game = NewTimedGame(30 * time.Second)
		m.state = StatePlaying
		return m, tickCmd()
	case "3":
		m.game = NewTimedGame(60 * time.Second)
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
		m.game = NewWordsGame(10)
		m.state = StatePlaying
		return m, tickCmd()
	case "2":
		m.game = NewWordsGame(25)
		m.state = StatePlaying
		return m, tickCmd()
	case "3":
		m.game = NewWordsGame(50)
		m.state = StatePlaying
		return m, tickCmd()
	case "4":
		m.game = NewWordsGame(100)
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
				if m.inputMode == "time" {
					m.game = NewTimedGame(time.Duration(value) * time.Second)
				} else {
					m.game = NewWordsGame(value)
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

func (m model) View() string {
	switch m.state {
	case StateMenu:
		return RenderMainMenu(m.leaderboard, m.width, m.height, m.wantToQuit)
	case StateTimeSelect:
		return RenderTimeSelect(m.leaderboard, m.width, m.height, m.wantToQuit)
	case StateWordsSelect:
		return RenderWordsSelect(m.leaderboard, m.width, m.height, m.wantToQuit)
	case StateCustomInput:
		return RenderCustomInput(m.customInput, m.inputMode, m.width, m.height)
	case StatePlaying:
		if m.game != nil {
			return RenderGame(m.game, m.width, m.height, m.wantToQuit)
		}
	case StateFinished:
		if m.game != nil {
			isPB := m.leaderboard.IsPB(m.game.WPM(), m.game.ModeString())
			return RenderFinished(m.game, m.width, m.height, isPB, m.wantToQuit)
		}
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
