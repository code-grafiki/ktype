package app

import (
	"ktype/internal/game"
	"ktype/internal/storage"
	"ktype/internal/ui"
)

// View returns the view for the current state
func (m Model) View() string {
	switch m.State {
	case game.StateMenu:
		return ui.RenderMainMenu(m.Leaderboard, m.Width, m.Height, m.WantToQuit, m.Difficulty, m.Complexity)
	case game.StateDifficultySelect:
		return ui.RenderDifficultySelect(m.Difficulty, m.Width, m.Height, m.WantToQuit)
	case game.StateComplexitySelect:
		return ui.RenderComplexitySelect(m.Complexity, m.Width, m.Height, m.WantToQuit)
	case game.StateStats:
		stats := storage.NewStatistics(m.Leaderboard)
		return ui.RenderStats(stats, m.Width, m.Height, m.WantToQuit)
	case game.StateHeatmap:
		return ui.RenderHeatmap(m.Heatmap, m.Width, m.Height, m.WantToQuit)
	case game.StateSettings:
		return ui.RenderSettings(m.ConfigManager, m.Width, m.Height, m.WantToQuit)
	case game.StateCursorSelect:
		return ui.RenderCursorSelect(m.ConfigManager, m.Width, m.Height, m.WantToQuit)
	case game.StateColorSelect:
		return ui.RenderColorSelect(m.ConfigManager, m.Width, m.Height, m.WantToQuit)
	case game.StateCustomWordList:
		return ui.RenderCustomWordList(m.WordListManager, m.CurrentWordList, m.Width, m.Height, m.WantToQuit)
	case game.StateTimeSelect:
		return ui.RenderTimeSelect(m.Leaderboard, m.Width, m.Height, m.WantToQuit)
	case game.StateWordsSelect:
		return ui.RenderWordsSelect(m.Leaderboard, m.Width, m.Height, m.WantToQuit)
	case game.StateCustomInput:
		return ui.RenderCustomInput(m.CustomInput, m.InputMode, m.Width, m.Height, m.ConfigManager.GetCursorType().CursorChar())
	case game.StatePlaying:
		if m.Game != nil {
			return ui.RenderGame(m.Game, m.Width, m.Height, m.WantToQuit, m.ConfigManager.GetCursorType().CursorChar())
		}
	case game.StateFinished:
		if m.Game != nil {
			isPB := m.Leaderboard.IsPB(m.Game.WPM(), m.Game.ModeString())
			return ui.RenderFinished(m.Game, m.Width, m.Height, isPB, m.WantToQuit)
		}
	case game.StateChallenges:
		return ui.RenderChallenges(m.Challenges, m.Width, m.Height, m.WantToQuit)
	}
	return ""
}
