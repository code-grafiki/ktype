package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// WordList represents a custom word list
type WordList struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Words       []string `json:"words"`
	CreatedAt   string   `json:"created_at"`
}

// WordListManager manages custom word lists
type WordListManager struct {
	Lists []WordList `json:"lists"`
	path  string
}

// NewWordListManager creates or loads a word list manager
func NewWordListManager() *WordListManager {
	wm := &WordListManager{
		Lists: []WordList{},
	}

	// Get config directory
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = os.TempDir()
	}

	ktypeDir := filepath.Join(configDir, "ktype")
	if err := os.MkdirAll(ktypeDir, 0755); err != nil {
		wm.path = "wordlists.json"
	} else {
		wm.path = filepath.Join(ktypeDir, "wordlists.json")
	}

	wm.load()
	return wm
}

// load reads word lists from file
func (wm *WordListManager) load() {
	data, err := os.ReadFile(wm.path)
	if err != nil {
		return // File doesn't exist yet
	}

	if err := json.Unmarshal(data, wm); err != nil {
		// Corrupt file - start fresh
		wm.Lists = []WordList{}
	}
}

// save writes word lists to file
func (wm *WordListManager) save() error {
	data, err := json.MarshalIndent(wm, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(wm.path, data, 0644)
}

// AddList adds a new word list
func (wm *WordListManager) AddList(name, description string, words []string) error {
	// Validate name
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("word list name cannot be empty")
	}

	// Check for duplicates
	for _, list := range wm.Lists {
		if list.Name == name {
			return fmt.Errorf("word list '%s' already exists", name)
		}
	}

	// Clean and validate words
	cleanWords := cleanWordList(words)
	if len(cleanWords) == 0 {
		return fmt.Errorf("word list must contain at least one word")
	}

	list := WordList{
		Name:        name,
		Description: description,
		Words:       cleanWords,
		CreatedAt:   fmt.Sprintf("%d", os.Getpid()), // Simple timestamp
	}

	wm.Lists = append(wm.Lists, list)
	return wm.save()
}

// DeleteList removes a word list by name
func (wm *WordListManager) DeleteList(name string) error {
	for i, list := range wm.Lists {
		if list.Name == name {
			wm.Lists = append(wm.Lists[:i], wm.Lists[i+1:]...)
			return wm.save()
		}
	}
	return fmt.Errorf("word list '%s' not found", name)
}

// GetList returns a word list by name
func (wm *WordListManager) GetList(name string) *WordList {
	for i := range wm.Lists {
		if wm.Lists[i].Name == name {
			return &wm.Lists[i]
		}
	}
	return nil
}

// ListNames returns all word list names
func (wm *WordListManager) ListNames() []string {
	names := make([]string, len(wm.Lists))
	for i, list := range wm.Lists {
		names[i] = list.Name
	}
	return names
}

// GetWords returns words from a list (or default words if not found)
func (wm *WordListManager) GetWords(name string, count int) []string {
	list := wm.GetList(name)
	if list == nil || len(list.Words) == 0 {
		return nil
	}

	// Return requested count (with repetition if needed)
	words := make([]string, count)
	for i := 0; i < count; i++ {
		words[i] = list.Words[i%len(list.Words)]
	}
	return words
}

// ImportFromFile imports a word list from a text file
func (wm *WordListManager) ImportFromFile(filepath, name, description string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word != "" && !strings.HasPrefix(word, "#") {
			words = append(words, word)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	return wm.AddList(name, description, words)
}

// ExportToFile exports a word list to a text file
func (wm *WordListManager) ExportToFile(name, filepath string) error {
	list := wm.GetList(name)
	if list == nil {
		return fmt.Errorf("word list '%s' not found", name)
	}

	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write header
	fmt.Fprintf(file, "# %s\n", list.Name)
	if list.Description != "" {
		fmt.Fprintf(file, "# %s\n", list.Description)
	}
	fmt.Fprintln(file, "#")

	// Write words
	for _, word := range list.Words {
		fmt.Fprintln(file, word)
	}

	return nil
}

// GetBuiltInLists returns the names of built-in word lists
func GetBuiltInLists() map[string][]string {
	return map[string][]string{
		"easy":   easyWords,
		"medium": mediumWords,
		"hard":   hardWords,
	}
}

// cleanWordList cleans and validates a word list
func cleanWordList(words []string) []string {
	var clean []string
	seen := make(map[string]bool)

	for _, word := range words {
		// Trim whitespace and lowercase
		word = strings.TrimSpace(strings.ToLower(word))

		// Skip empty words
		if word == "" {
			continue
		}

		// Skip duplicates
		if seen[word] {
			continue
		}

		// Remove any internal whitespace
		word = strings.Join(strings.Fields(word), "")

		clean = append(clean, word)
		seen[word] = true
	}

	return clean
}

// ValidateWordListName checks if a name is valid
func ValidateWordListName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("name cannot be empty")
	}

	if len(name) > 50 {
		return fmt.Errorf("name too long (max 50 characters)")
	}

	// Check for invalid characters
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range invalidChars {
		if strings.Contains(name, char) {
			return fmt.Errorf("name contains invalid character '%s'", char)
		}
	}

	return nil
}

// Count returns the number of word lists
func (wm *WordListManager) Count() int {
	return len(wm.Lists)
}

// Clear removes all custom word lists
func (wm *WordListManager) Clear() error {
	wm.Lists = []WordList{}
	return wm.save()
}
