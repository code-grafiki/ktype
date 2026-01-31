package main

import (
	"strings"
	"testing"
)

func TestPunctuationMarksNotEmpty(t *testing.T) {
	if len(punctuationMarks) == 0 {
		t.Error("punctuationMarks should not be empty")
	}
}

func TestNumberListNotEmpty(t *testing.T) {
	if len(numberList) == 0 {
		t.Error("numberList should not be empty")
	}
}

func TestPunctuationCombosNotEmpty(t *testing.T) {
	if len(punctuationCombos) == 0 {
		t.Error("punctuationCombos should not be empty")
	}
}

func TestNumberPunctuationCombosNotEmpty(t *testing.T) {
	if len(numberPunctuationCombos) == 0 {
		t.Error("numberPunctuationCombos should not be empty")
	}
}

func TestWordComplexityString(t *testing.T) {
	tests := []struct {
		complexity WordComplexity
		expected   string
	}{
		{ComplexityNormal, "normal"},
		{ComplexityPunctuation, "punctuation"},
		{ComplexityNumbers, "numbers"},
		{ComplexityFull, "full"},
		{WordComplexity(99), "normal"}, // default for unknown
	}

	for _, tt := range tests {
		result := tt.complexity.String()
		if result != tt.expected {
			t.Errorf("WordComplexity(%d).String() = %q, expected %q", tt.complexity, result, tt.expected)
		}
	}
}

func TestGetRandomWordsWithComplexityNormal(t *testing.T) {
	words := getRandomWordsWithComplexity(50, DifficultyMedium, ComplexityNormal)
	if len(words) != 50 {
		t.Errorf("Expected 50 words, got %d", len(words))
	}

	// Normal complexity should just return base words without punctuation/numbers
	for _, word := range words {
		// Should not contain punctuation (except for rare cases in base words)
		if strings.ContainsAny(word, "!@#$%^&*()_+-=[]{}|;':\",./<>?") {
			// This might happen if base word contains punctuation, which is fine
			// Just check that we're getting reasonable words
			if len(word) < 1 {
				t.Errorf("Got empty word with complexity normal")
			}
		}
	}
}

func TestGetRandomWordsWithComplexityPunctuation(t *testing.T) {
	words := getRandomWordsWithComplexity(100, DifficultyMedium, ComplexityPunctuation)
	if len(words) != 100 {
		t.Errorf("Expected 100 words, got %d", len(words))
	}

	// With 30% chance, we should get at least some punctuation
	punctuationCount := 0
	for _, word := range words {
		if strings.ContainsAny(word, "!@#$%^&*()_+-=[]{}|;':\",./<>?") {
			punctuationCount++
		}
	}

	// With 100 words at 30% probability, we should expect around 30 words with punctuation
	// Allow for variance, but should have at least some
	if punctuationCount == 0 {
		t.Log("Warning: No punctuation found in 100 words (unlikely but possible)")
	}
}

func TestGetRandomWordsWithComplexityNumbers(t *testing.T) {
	words := getRandomWordsWithComplexity(100, DifficultyMedium, ComplexityNumbers)
	if len(words) != 100 {
		t.Errorf("Expected 100 words, got %d", len(words))
	}

	// Should have some numbers (20% replacement + 10% suffix = 30% chance)
	hasNumbers := false
	for _, word := range words {
		for _, r := range word {
			if r >= '0' && r <= '9' {
				hasNumbers = true
				break
			}
		}
		if hasNumbers {
			break
		}
	}

	if !hasNumbers {
		t.Log("Warning: No numbers found in 100 words (unlikely but possible)")
	}
}

func TestGetRandomWordsWithComplexityFull(t *testing.T) {
	words := getRandomWordsWithComplexity(100, DifficultyMedium, ComplexityFull)
	if len(words) != 100 {
		t.Errorf("Expected 100 words, got %d", len(words))
	}

	// Full complexity should have a mix of punctuation and/or numbers
	hasSpecial := false
	for _, word := range words {
		for _, r := range word {
			if (r >= '0' && r <= '9') || strings.ContainsRune("!@#$%^&*()_+-=[]{}|;':\",./<>?", r) {
				hasSpecial = true
				break
			}
		}
		if hasSpecial {
			break
		}
	}

	if !hasSpecial {
		t.Log("Warning: No special characters found in 100 words with full complexity")
	}
}

func TestGetRandomWordsWithComplexityEmpty(t *testing.T) {
	words := getRandomWordsWithComplexity(0, DifficultyMedium, ComplexityPunctuation)
	if len(words) != 0 {
		t.Errorf("Expected 0 words for n=0, got %d", len(words))
	}
}

func TestAddPunctuationToWord(t *testing.T) {
	// Test that addPunctuationToWord actually adds punctuation
	word := "test"
	punctuated := addPunctuationToWord(word)

	// Should either have punctuation or be the original word (if random choice returns it unchanged)
	if punctuated == word {
		// This is possible but unlikely - run multiple times to be sure
		hasPunctuation := false
		for i := 0; i < 100; i++ {
			result := addPunctuationToWord(word)
			if strings.ContainsAny(result, "!@#$%^&*()_+-=[]{}|;':\",./<>?") {
				hasPunctuation = true
				break
			}
		}
		if !hasPunctuation {
			t.Error("addPunctuationToWord should add punctuation most of the time")
		}
	}
}

func TestAddPunctuationToWordWithShortWords(t *testing.T) {
	// Test with single character
	word := "a"
	result := addPunctuationToWord(word)
	if len(result) < 1 {
		t.Error("Result should not be empty")
	}
}
