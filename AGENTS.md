# AGENTS.md - Development Guidelines for ktype

## Project Overview

ktype is a CLI typing test application built with Go using the Bubble Tea TUI framework. It features multiple game modes (timed, words, zen), a leaderboard system with personal bests, statistics tracking, heatmaps, daily challenges, and a clean MonkeyType-inspired interface.

## Build Commands

```bash
# Build the application
go build -o ktype .

# Run the application
go run .

# Build for release (with optimizations)
go build -ldflags "-s -w" -o ktype .
```

## Test Commands

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test by name
go test -run TestFunctionName ./...

# Run specific test in a specific file
go test -run TestFunctionName path/to/file_test.go

# Run tests in verbose mode
go test -v ./...

# Run benchmarks
go test -bench=. ./...
```

## Lint Commands

```bash
# Format code
go fmt ./...

# Vet code for suspicious constructs
go vet ./...

# Run golangci-lint (if installed)
golangci-lint run

# Tidy go modules
go mod tidy

# Verify dependencies
go mod verify
```

## Code Style Guidelines

### General Go Conventions

- Use `gofmt` for automatic formatting
- Follow standard Go naming conventions
- Keep functions focused and under 50 lines when possible
- Maximum line length: ~100 characters (soft limit)

### Naming Conventions

- **Exported types/functions**: PascalCase (e.g., `NewGame`, `GameState`)
- **Unexported types/functions**: camelCase (e.g., `handleKeyPress`, `wordIndex`)
- **Constants**: PascalCase for exported, camelCase for unexported
- **Acronyms**: All caps (e.g., `WPM`, `HTTP`, `ID`)
- **Interfaces**: Descriptive noun or -er suffix (e.g., `Reader`, `ScoreManager`)

### Imports Organization

```go
import (
    // Standard library packages (alphabetically)
    "fmt"
    "os"
    "time"

    // Third-party packages (alphabetically)
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)
```

### Types and Structs

- Define types at the top of files after imports
- Add comments for exported types
- Group related constants with `iota`
- Document struct fields inline when not obvious

```go
// GameState represents the current state of the game
type GameState int

const (
    StateMenu GameState = iota
    StatePlaying
    StateFinished
)
```

### Error Handling

- Check errors immediately after assignment
- Return errors rather than logging when possible
- Use descriptive error messages
- Handle file operations with proper fallbacks

```go
// Good
configDir, err := os.UserConfigDir()
if err != nil {
    configDir = os.TempDir()
}

// Avoid panics - handle gracefully
if err := os.MkdirAll(ktypeDir, 0755); err != nil {
    lb.path = "scores.json"
}
```

### Comments

- Use complete sentences with periods
- Exported items must have comments
- Comments should describe purpose, not implementation
- Keep comments current with code changes

### Bubble Tea Patterns

- Model struct contains all state
- Update method handles all message types
- Use type switch for message handling
- Keep Init simple, return nil if no initial commands
- View should be pure function of model state

### UI/Styling (lipgloss)

- Define color palette as package-level variables
- Group related styles together
- Use semantic names (e.g., `errorStyle`, `accentStyle`)
- Chain style modifications for variants

```go
var (
    colorError   = lipgloss.Color("#ca4754")
    errorStyle   = lipgloss.NewStyle().Foreground(colorError)
    warningStyle = errorStyle.Copy().Bold(true)
)
```

## File Organization

- `main.go`: Application entry point, Bubble Tea model
- `game.go`: Game logic, state management, scoring
- `ui.go`: Rendering functions, lipgloss styles
- `leaderboard.go`: Score persistence, personal bests
- `words.go`: Word lists and random selection
- `wordlist.go`: Custom word list management
- `statistics.go`: Statistics tracking and analytics
- `heatmap.go`: Typing heatmap visualization
- `config.go`: Configuration management
- `challenges.go`: Daily challenges system
- `*_test.go`: Test files (create as needed)

## Dependencies

Key dependencies (managed via go.mod):
- `github.com/charmbracelet/bubbletea`: TUI framework
- `github.com/charmbracelet/lipgloss`: Styling library
- `github.com/charmbracelet/colorprofile`: Color profile support

## Git Workflow

```bash
# Before committing
go fmt ./...
go vet ./...
go test ./...

# Commit messages: concise, imperative mood
# Good: "Add zen mode game type"
# Bad: "Added zen mode" or "I added zen mode"
```

## Testing Guidelines

- Test exported functions primarily
- Use table-driven tests for multiple cases
- Mock external dependencies (file system, time)
- Test edge cases and error conditions
- Name tests descriptively: `TestFunctionName_Scenario`
