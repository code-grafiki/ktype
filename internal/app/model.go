package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"ktype/internal/game"
	"ktype/internal/storage"
	"ktype/internal/words"
)

// tickMsg is sent every tick to update the timer
type tickMsg time.Time

// Model is the Bubble Tea model
type Model struct {
	Game        *game.Game
	Leaderboard *storage.Leaderboard
	Width       int
	Height      int
	State       game.State
	WantToQuit  bool
	QuitPressAt time.Time

	// For custom input
	CustomInput string
	InputMode   string
	Difficulty  words.Difficulty
	Complexity  words.Complexity

	// For custom word lists
	WordListManager *storage.WordListManager
	CurrentWordList string

	// For heatmap
	Heatmap *storage.Heatmap

	// For configuration
	ConfigManager *storage.ConfigManager

	// Daily challenges
	Challenges *storage.DailyChallenges
}

// InitialModel creates the initial model
func InitialModel() Model {
	cm := storage.NewConfigManager()

	return Model{
		State:           game.StateMenu,
		Width:           80,
		Height:          24,
		Leaderboard:     storage.NewLeaderboard(),
		Difficulty:      words.DifficultyMedium,
		Complexity:      words.ComplexityNormal,
		WordListManager: storage.NewWordListManager(),
		CurrentWordList: "",
		Heatmap:         storage.NewHeatmap(),
		ConfigManager:   cm,
		Challenges:      storage.NewDailyChallenges(),
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}
