package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// CursorType represents different cursor styles
type CursorType int

const (
	CursorBlock CursorType = iota
	CursorUnderscore
	CursorLine
	CursorBeam
	CursorUnderline
	CursorBar
)

// String returns the cursor type name
func (c CursorType) String() string {
	switch c {
	case CursorBlock:
		return "block"
	case CursorUnderscore:
		return "underscore"
	case CursorLine:
		return "line"
	case CursorBeam:
		return "beam"
	case CursorUnderline:
		return "underline"
	case CursorBar:
		return "bar"
	default:
		return "underscore"
	}
}

// CursorChar returns the character to use for the cursor
func (c CursorType) CursorChar() string {
	switch c {
	case CursorBlock:
		return "█"
	case CursorUnderscore:
		return "_"
	case CursorLine:
		return "|"
	case CursorBeam:
		return "┃"
	case CursorUnderline:
		return "_"
	case CursorBar:
		return "|"
	default:
		return "_"
	}
}

// AccentColor represents preset accent colors
type AccentColor int

const (
	ColorRed AccentColor = iota
	ColorOrange
	ColorYellow
	ColorGreen
	ColorCyan
	ColorBlue
	ColorPurple
	ColorPink
	ColorWhite
	ColorBlack
	ColorCustom
)

// Config holds user preferences
type Config struct {
	CursorType      CursorType  `json:"cursor_type"`
	AccentColor     string      `json:"accent_color"`
	AccentColorEnum AccentColor `json:"accent_color_enum"`
	CustomColor     string      `json:"custom_color,omitempty"`
	ShowHeatmap     bool        `json:"show_heatmap"`
	SoundEnabled    bool        `json:"sound_enabled"`
}

// DefaultConfig returns default configuration
func DefaultConfig() Config {
	return Config{
		CursorType:      CursorUnderscore,
		AccentColor:     "#5eacd3",
		AccentColorEnum: ColorCyan,
		CustomColor:     "",
		ShowHeatmap:     true,
		SoundEnabled:    false,
	}
}

// ConfigManager manages user configuration
type ConfigManager struct {
	config Config
	path   string
}

// NewConfigManager creates or loads configuration
func NewConfigManager() *ConfigManager {
	cm := &ConfigManager{
		config: DefaultConfig(),
	}

	// Get config directory
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = os.TempDir()
	}

	ktypeDir := filepath.Join(configDir, "ktype")
	if err := os.MkdirAll(ktypeDir, 0755); err != nil {
		cm.path = "config.json"
	} else {
		cm.path = filepath.Join(ktypeDir, "config.json")
	}

	cm.load()
	return cm
}

// load reads config from file
func (cm *ConfigManager) load() {
	data, err := os.ReadFile(cm.path)
	if err != nil {
		return // File doesn't exist yet, use defaults
	}

	if err := json.Unmarshal(data, &cm.config); err != nil {
		// Corrupt file - use defaults
		cm.config = DefaultConfig()
	}
}

// save writes config to file
func (cm *ConfigManager) save() error {
	data, err := json.MarshalIndent(cm.config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cm.path, data, 0644)
}

// GetConfig returns the current configuration
func (cm *ConfigManager) GetConfig() Config {
	return cm.config
}

// SetCursorType updates the cursor type
func (cm *ConfigManager) SetCursorType(ct CursorType) error {
	cm.config.CursorType = ct
	return cm.save()
}

// SetAccentColor updates the accent color (string version)
func (cm *ConfigManager) SetAccentColorString(color string) error {
	cm.config.AccentColor = color
	return cm.save()
}

// SetAccentColor updates the accent color (enum version)
func (cm *ConfigManager) SetAccentColor(color AccentColor) {
	cm.config.AccentColorEnum = color
	if color == ColorCustom && cm.config.CustomColor != "" {
		cm.config.AccentColor = cm.config.CustomColor
	} else {
		// Set default color for presets
		colors := []string{
			"#e06c75", // Red
			"#d19a66", // Orange
			"#e5c07b", // Yellow
			"#98c379", // Green
			"#56b6c2", // Cyan
			"#61afef", // Blue
			"#c678dd", // Purple
			"#ff79c6", // Pink
			"#abb2bf", // White
			"#282c34", // Black
		}
		if int(color) >= 0 && int(color) < len(colors) {
			cm.config.AccentColor = colors[color]
		}
	}
	cm.save()
}

// SetCustomColor sets a custom hex color, returns true if valid
func (cm *ConfigManager) SetCustomColor(hex string) bool {
	if ValidateColor(hex) {
		cm.config.CustomColor = hex
		cm.save()
		return true
	}
	return false
}

// GetCursorType returns the current cursor type
func (cm *ConfigManager) GetCursorType() CursorType {
	return cm.config.CursorType
}

// GetAccentColor returns the current accent color (string version)
func (cm *ConfigManager) GetAccentColor() string {
	return cm.config.AccentColor
}

// GetAccentColorEnum returns the current accent color (enum version)
func (cm *ConfigManager) GetAccentColorEnum() AccentColor {
	return cm.config.AccentColorEnum
}

// GetCursorTypeName returns the display name for the current cursor type
func (cm *ConfigManager) GetCursorTypeName() string {
	names := []string{"Block", "Underscore", "Line", "Beam", "Underline", "Bar"}
	if cm.config.CursorType >= 0 && int(cm.config.CursorType) < len(names) {
		return names[cm.config.CursorType]
	}
	return "Underscore"
}

// GetAccentColorName returns the display name for the current accent color
func (cm *ConfigManager) GetAccentColorName() string {
	names := []string{"Red", "Orange", "Yellow", "Green", "Cyan", "Blue", "Purple", "Pink", "White", "Black", "Custom"}
	if cm.config.AccentColorEnum >= 0 && int(cm.config.AccentColorEnum) < len(names) {
		return names[cm.config.AccentColorEnum]
	}
	return "Cyan"
}

// GetAccentColorHex returns the hex code for the current accent color
func (cm *ConfigManager) GetAccentColorHex() string {
	if cm.config.AccentColorEnum == ColorCustom && cm.config.CustomColor != "" {
		return cm.config.CustomColor
	}

	colors := []string{
		"#e06c75", // Red
		"#d19a66", // Orange
		"#e5c07b", // Yellow
		"#98c379", // Green
		"#56b6c2", // Cyan
		"#61afef", // Blue
		"#c678dd", // Purple
		"#ff79c6", // Pink
		"#abb2bf", // White
		"#282c34", // Black
	}

	if cm.config.AccentColorEnum >= 0 && int(cm.config.AccentColorEnum) < len(colors) {
		return colors[cm.config.AccentColorEnum]
	}
	return "#56b6c2" // Default cyan
}

// ValidateColor checks if a color string is valid hex
func ValidateColor(color string) bool {
	if len(color) != 7 && len(color) != 4 {
		return false
	}
	if color[0] != '#' {
		return false
	}
	for i := 1; i < len(color); i++ {
		c := color[i]
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

// GetCursorTypeFromString parses cursor type from string
func GetCursorTypeFromString(s string) CursorType {
	switch s {
	case "block":
		return CursorBlock
	case "line":
		return CursorLine
	case "beam":
		return CursorBeam
	default:
		return CursorUnderscore
	}
}

// PresetColors returns a list of preset accent colors
func PresetColors() []string {
	return []string{
		"#5eacd3", // Blue (default)
		"#98c379", // Green
		"#e2b714", // Yellow
		"#d19a66", // Orange
		"#ca4754", // Red
		"#c678dd", // Purple
		"#56b6c2", // Cyan
		"#ffffff", // White
		"#ff6b6b", // Coral
		"#4ecdc4", // Turquoise
	}
}
