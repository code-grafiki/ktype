package words

import (
	"strings"
	"testing"
)

func TestComplexityString(t *testing.T) {
	tests := []struct {
		complexity Complexity
		expected   string
	}{
		{ComplexityNormal, "normal"},
		{ComplexityPunctuation, "punctuation"},
		{ComplexityNumbers, "numbers"},
		{ComplexityFull, "full"},
		{Complexity(99), "normal"}, // default for unknown
	}

	for _, tt := range tests {
		result := tt.complexity.String()
		if result != tt.expected {
			t.Errorf("Complexity(%d).String() = %q, expected %q", tt.complexity, result, tt.expected)
		}
	}
}

func TestGetRandomWithComplexityNormal(t *testing.T) {
	words := GetRandomWithComplexity(50, DifficultyMedium, ComplexityNormal)
	if len(words) != 50 {
		t.Errorf("Expected 50 words, got %d", len(words))
	}

	// Normal complexity should just return base words without punctuation/numbers
	for _, word := range words {
		if len(word) < 1 {
			t.Errorf("Got empty word with complexity normal")
		}
	}
}

func TestGetRandomWithComplexityPunctuation(t *testing.T) {
	words := GetRandomWithComplexity(100, DifficultyMedium, ComplexityPunctuation)
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

	if punctuationCount == 0 {
		t.Log("Warning: No punctuation found in 100 words (unlikely but possible)")
	}
}

func TestGetRandomWithComplexityNumbers(t *testing.T) {
	words := GetRandomWithComplexity(100, DifficultyMedium, ComplexityNumbers)
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

func TestGetRandomWithComplexityFull(t *testing.T) {
	words := GetRandomWithComplexity(100, DifficultyMedium, ComplexityFull)
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

func TestGetRandomWithComplexityEmpty(t *testing.T) {
	words := GetRandomWithComplexity(0, DifficultyMedium, ComplexityPunctuation)
	if len(words) != 0 {
		t.Errorf("Expected 0 words for n=0, got %d", len(words))
	}
}

func TestAddPunctuation(t *testing.T) {
	// Test that AddPunctuation actually adds punctuation
	word := "test"
	punctuated := AddPunctuation(word)

	// Should either have punctuation or be the original word (if random choice returns it unchanged)
	if punctuated == word {
		// This is possible but unlikely - run multiple times to be sure
		hasPunctuation := false
		for i := 0; i < 100; i++ {
			result := AddPunctuation(word)
			if strings.ContainsAny(result, "!@#$%^&*()_+-=[]{}|;':\",./<>?") {
				hasPunctuation = true
				break
			}
		}
		if !hasPunctuation {
			t.Error("AddPunctuation should add punctuation most of the time")
		}
	}
}
