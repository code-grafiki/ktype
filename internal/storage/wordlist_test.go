package storage

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewWordListManager(t *testing.T) {
	wm := NewWordListManager()
	if wm == nil {
		t.Fatal("NewWordListManager returned nil")
	}
	if wm.Lists == nil {
		t.Error("WordListManager.Lists should be initialized")
	}
}

func TestWordListManagerAddList(t *testing.T) {
	tempDir := t.TempDir()
	wm := &WordListManager{
		Lists: []WordList{},
		path:  filepath.Join(tempDir, "wordlists.json"),
	}

	words := []string{"apple", "banana", "cherry"}
	err := wm.AddList("fruits", "common fruits", words)
	if err != nil {
		t.Fatalf("Failed to add list: %v", err)
	}

	if wm.Count() != 1 {
		t.Errorf("Expected 1 list, got %d", wm.Count())
	}

	list := wm.GetList("fruits")
	if list == nil {
		t.Fatal("Could not retrieve added list")
	}

	if list.Name != "fruits" {
		t.Errorf("Expected name 'fruits', got %s", list.Name)
	}

	if len(list.Words) != 3 {
		t.Errorf("Expected 3 words, got %d", len(list.Words))
	}
}

func TestWordListManagerAddListDuplicate(t *testing.T) {
	tempDir := t.TempDir()
	wm := &WordListManager{
		Lists: []WordList{},
		path:  filepath.Join(tempDir, "wordlists.json"),
	}

	words := []string{"apple", "banana"}
	wm.AddList("fruits", "common fruits", words)

	// Try to add duplicate
	err := wm.AddList("fruits", "duplicate", words)
	if err == nil {
		t.Error("Should fail when adding duplicate list name")
	}
}

func TestWordListManagerAddListEmptyName(t *testing.T) {
	tempDir := t.TempDir()
	wm := &WordListManager{
		Lists: []WordList{},
		path:  filepath.Join(tempDir, "wordlists.json"),
	}

	words := []string{"apple", "banana"}
	err := wm.AddList("", "no name", words)
	if err == nil {
		t.Error("Should fail when adding list with empty name")
	}
}

func TestWordListManagerAddListEmptyWords(t *testing.T) {
	tempDir := t.TempDir()
	wm := &WordListManager{
		Lists: []WordList{},
		path:  filepath.Join(tempDir, "wordlists.json"),
	}

	err := wm.AddList("empty", "no words", []string{})
	if err == nil {
		t.Error("Should fail when adding list with no words")
	}
}

func TestWordListManagerDeleteList(t *testing.T) {
	tempDir := t.TempDir()
	wm := &WordListManager{
		Lists: []WordList{},
		path:  filepath.Join(tempDir, "wordlists.json"),
	}

	words := []string{"apple", "banana"}
	wm.AddList("fruits", "common fruits", words)

	err := wm.DeleteList("fruits")
	if err != nil {
		t.Fatalf("Failed to delete list: %v", err)
	}

	if wm.Count() != 0 {
		t.Errorf("Expected 0 lists after deletion, got %d", wm.Count())
	}

	// Try to delete non-existent
	err = wm.DeleteList("nonexistent")
	if err == nil {
		t.Error("Should fail when deleting non-existent list")
	}
}

func TestWordListManagerListNames(t *testing.T) {
	tempDir := t.TempDir()
	wm := &WordListManager{
		Lists: []WordList{},
		path:  filepath.Join(tempDir, "wordlists.json"),
	}

	wm.AddList("fruits", "common fruits", []string{"apple", "banana"})
	wm.AddList("colors", "basic colors", []string{"red", "blue"})

	names := wm.ListNames()
	if len(names) != 2 {
		t.Errorf("Expected 2 names, got %d", len(names))
	}

	// Check that both names are present
	foundFruits := false
	foundColors := false
	for _, name := range names {
		if name == "fruits" {
			foundFruits = true
		}
		if name == "colors" {
			foundColors = true
		}
	}

	if !foundFruits || !foundColors {
		t.Error("Expected to find both 'fruits' and 'colors' in list names")
	}
}

func TestWordListManagerGetWords(t *testing.T) {
	tempDir := t.TempDir()
	wm := &WordListManager{
		Lists: []WordList{},
		path:  filepath.Join(tempDir, "wordlists.json"),
	}

	words := []string{"apple", "banana", "cherry"}
	wm.AddList("fruits", "common fruits", words)

	// Get more words than in list (should cycle)
	result := wm.GetWords("fruits", 5)
	if len(result) != 5 {
		t.Errorf("Expected 5 words, got %d", len(result))
	}

	// Check cycling works
	if result[0] != "apple" || result[3] != "apple" {
		t.Error("Words should cycle when requesting more than available")
	}

	// Try non-existent list
	result = wm.GetWords("nonexistent", 3)
	if result != nil {
		t.Error("Should return nil for non-existent list")
	}
}

func TestWordListManagerSaveAndLoad(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "wordlists.json")

	// Create and populate
	wm1 := &WordListManager{
		Lists: []WordList{},
		path:  path,
	}

	wm1.AddList("fruits", "common fruits", []string{"apple", "banana"})
	err := wm1.save()
	if err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	// Load into new manager
	wm2 := &WordListManager{
		Lists: []WordList{},
		path:  path,
	}
	wm2.load()

	if wm2.Count() != 1 {
		t.Errorf("Expected 1 list after load, got %d", wm2.Count())
	}

	list := wm2.GetList("fruits")
	if list == nil {
		t.Fatal("Could not retrieve loaded list")
	}

	if len(list.Words) != 2 {
		t.Errorf("Expected 2 words after load, got %d", len(list.Words))
	}
}

func TestCleanWordList(t *testing.T) {
	input := []string{
		"  Apple  ",
		"  apple  ", // duplicate
		"Banana",
		"", // empty
		"cherry",
		"# comment", // should be kept as a word (not a file line)
	}

	result := cleanWordList(input)

	if len(result) != 4 {
		t.Errorf("Expected 4 words after cleaning, got %d", len(result))
	}

	// Check all words are lowercase
	for _, word := range result {
		if word != strings.ToLower(word) {
			t.Errorf("Word '%s' should be lowercase", word)
		}
	}
}

func TestValidateWordListName(t *testing.T) {
	tests := []struct {
		name    string
		isValid bool
	}{
		{"valid", true},
		{"", false},                      // empty
		{"a/b", false},                   // path separator
		{"a\\b", false},                  // backslash
		{"a:b", false},                   // colon
		{"a*b", false},                   // asterisk
		{"a?b", false},                   // question mark
		{"a\"b", false},                  // quote
		{"a<b", false},                   // less than
		{"a>b", false},                   // greater than
		{"a|b", false},                   // pipe
		{strings.Repeat("a", 51), false}, // too long
	}

	for _, tt := range tests {
		err := ValidateWordListName(tt.name)
		if tt.isValid && err != nil {
			t.Errorf("Name '%s' should be valid, got error: %v", tt.name, err)
		}
		if !tt.isValid && err == nil {
			t.Errorf("Name '%s' should be invalid, but no error", tt.name)
		}
	}
}

func TestWordListManagerClear(t *testing.T) {
	tempDir := t.TempDir()
	wm := &WordListManager{
		Lists: []WordList{
			{Name: "fruits", Words: []string{"apple"}},
			{Name: "colors", Words: []string{"red"}},
		},
		path: filepath.Join(tempDir, "wordlists.json"),
	}

	err := wm.Clear()
	if err != nil {
		t.Fatalf("Failed to clear: %v", err)
	}

	if wm.Count() != 0 {
		t.Errorf("Expected 0 lists after clear, got %d", wm.Count())
	}
}

func TestWordListManagerImportFromFile(t *testing.T) {
	tempDir := t.TempDir()
	wm := &WordListManager{
		Lists: []WordList{},
		path:  filepath.Join(tempDir, "wordlists.json"),
	}

	// Create test file
	testFile := filepath.Join(tempDir, "test_words.txt")
	content := "# Fruits\napple\nbanana\ncherry\n"
	os.WriteFile(testFile, []byte(content), 0644)

	err := wm.ImportFromFile(testFile, "imported", "imported from file")
	if err != nil {
		t.Fatalf("Failed to import: %v", err)
	}

	list := wm.GetList("imported")
	if list == nil {
		t.Fatal("Could not retrieve imported list")
	}

	if len(list.Words) != 3 {
		t.Errorf("Expected 3 words, got %d", len(list.Words))
	}
}

func TestWordListManagerExportToFile(t *testing.T) {
	tempDir := t.TempDir()
	wm := &WordListManager{
		Lists: []WordList{},
		path:  filepath.Join(tempDir, "wordlists.json"),
	}

	wm.AddList("fruits", "common fruits", []string{"apple", "banana"})

	exportFile := filepath.Join(tempDir, "exported.txt")
	err := wm.ExportToFile("fruits", exportFile)
	if err != nil {
		t.Fatalf("Failed to export: %v", err)
	}

	// Check file was created and contains content
	content, err := os.ReadFile(exportFile)
	if err != nil {
		t.Fatalf("Could not read exported file: %v", err)
	}

	if !strings.Contains(string(content), "apple") {
		t.Error("Exported file should contain 'apple'")
	}
}
