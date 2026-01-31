package words

import (
	"testing"
)

func TestGetRandom(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		expected int
	}{
		{"zero words", 0, 0},
		{"one word", 1, 1},
		{"ten words", 10, 10},
		{"fifty words", 50, 50},
		{"hundred words", 100, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			words := GetRandom(tt.n, DifficultyMedium)

			if len(words) != tt.expected {
				t.Errorf("GetRandom(%d) returned %d words, expected %d",
					tt.n, len(words), tt.expected)
			}
		})
	}
}

func TestGetRandomContent(t *testing.T) {
	// Test all difficulty levels
	difficulties := []Difficulty{DifficultyEasy, DifficultyMedium, DifficultyHard}

	for _, difficulty := range difficulties {
		words := GetRandom(50, difficulty)
		wordSet := make(map[string]bool)
		list := GetList(difficulty)
		for _, word := range list {
			wordSet[word] = true
		}

		for _, word := range words {
			if !wordSet[word] {
				t.Errorf("Word %q is not in the %s word list", word, difficulty.String())
			}
		}
	}
}

func TestGetRandomRandomness(t *testing.T) {
	// Get two sets of random words
	words1 := GetRandom(20, DifficultyMedium)
	words2 := GetRandom(20, DifficultyMedium)

	// They should be different (with very high probability)
	allSame := true
	for i := 0; i < len(words1) && i < len(words2); i++ {
		if words1[i] != words2[i] {
			allSame = false
			break
		}
	}

	if allSame {
		t.Error("Two calls to GetRandom returned identical results - randomness may not be working")
	}
}

func TestGetRandomNoConsecutiveDuplicates(t *testing.T) {
	// Generate a large number of words to check for consecutive duplicates
	words := GetRandom(200, DifficultyMedium)

	consecutiveDuplicates := 0
	for i := 1; i < len(words); i++ {
		if words[i] == words[i-1] {
			consecutiveDuplicates++
			t.Errorf("Consecutive duplicate found at positions %d and %d: %q",
				i-1, i, words[i])
		}
	}

	if consecutiveDuplicates > 0 {
		t.Errorf("Found %d consecutive duplicate pairs", consecutiveDuplicates)
	}
}

func TestWordListsNotEmpty(t *testing.T) {
	lists := []struct {
		name       string
		difficulty Difficulty
	}{
		{"easy", DifficultyEasy},
		{"medium", DifficultyMedium},
		{"hard", DifficultyHard},
	}

	for _, list := range lists {
		words := GetList(list.difficulty)
		if len(words) == 0 {
			t.Errorf("%s word list is empty", list.name)
		}

		if len(words) < 50 {
			t.Errorf("%s word list has only %d words, expected at least 50", list.name, len(words))
		}
	}
}

func TestWordListsNoEmptyStrings(t *testing.T) {
	difficulties := []Difficulty{DifficultyEasy, DifficultyMedium, DifficultyHard}
	listNames := []string{"easy", "medium", "hard"}

	for idx, difficulty := range difficulties {
		wordList := GetList(difficulty)
		for i, word := range wordList {
			if word == "" {
				t.Errorf("%s wordList[%d] is empty string", listNames[idx], i)
			}
		}
	}
}

func TestWordListsNoDuplicates(t *testing.T) {
	difficulties := []Difficulty{DifficultyEasy, DifficultyMedium, DifficultyHard}
	listNames := []string{"easy", "medium", "hard"}

	for idx, difficulty := range difficulties {
		wordList := GetList(difficulty)
		seen := make(map[string]int)
		duplicates := []string{}

		for i, word := range wordList {
			if prevIdx, exists := seen[word]; exists {
				duplicates = append(duplicates, word)
				t.Errorf("Duplicate word %q found in %s list at indices %d and %d", word, listNames[idx], prevIdx, i)
			}
			seen[word] = i
		}

		if len(duplicates) > 0 {
			t.Errorf("Found %d duplicate words in %s list: %v", len(duplicates), listNames[idx], duplicates)
		}
	}
}

func TestWordListsMostlyLowercase(t *testing.T) {
	difficulties := []Difficulty{DifficultyEasy, DifficultyMedium, DifficultyHard}
	listNames := []string{"easy", "medium", "hard"}

	for idx, difficulty := range difficulties {
		wordList := GetList(difficulty)
		// Count words with uppercase letters (only "I" should be allowed)
		uppercaseCount := 0
		for i, word := range wordList {
			for _, r := range word {
				if r >= 'A' && r <= 'Z' {
					uppercaseCount++
					// Only "I" is allowed to be uppercase (pronoun)
					if word != "I" {
						t.Errorf("%s wordList[%d] = %q contains unexpected uppercase letter", listNames[idx], i, word)
					}
					break
				}
			}
		}
		// Should have exactly one uppercase word: "I" (only in easy list)
		if listNames[idx] == "easy" && uppercaseCount > 1 {
			t.Errorf("Expected at most 1 uppercase word (I) in easy list, found %d", uppercaseCount)
		}
	}
}

func TestWordListsNoWhitespace(t *testing.T) {
	difficulties := []Difficulty{DifficultyEasy, DifficultyMedium, DifficultyHard}
	listNames := []string{"easy", "medium", "hard"}

	for idx, difficulty := range difficulties {
		wordList := GetList(difficulty)
		for i, word := range wordList {
			for _, r := range word {
				if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
					t.Errorf("%s wordList[%d] = %q contains whitespace", listNames[idx], i, word)
					break
				}
			}
		}
	}
}

func TestGetList(t *testing.T) {
	// Check that GetList returns the correct lists by checking lengths
	easy := GetList(DifficultyEasy)
	medium := GetList(DifficultyMedium)
	hard := GetList(DifficultyHard)

	if len(easy) == 0 {
		t.Error("GetList(DifficultyEasy) returned empty list")
	}
	if len(medium) == 0 {
		t.Error("GetList(DifficultyMedium) returned empty list")
	}
	if len(hard) == 0 {
		t.Error("GetList(DifficultyHard) returned empty list")
	}

	// Check they return different lists (different lengths or first elements)
	if len(easy) == len(medium) && len(easy) > 0 && easy[0] == medium[0] {
		t.Error("GetList should return different lists for different difficulties")
	}
}

func TestDifficultyString(t *testing.T) {
	tests := []struct {
		difficulty Difficulty
		expected   string
	}{
		{DifficultyEasy, "easy"},
		{DifficultyMedium, "medium"},
		{DifficultyHard, "hard"},
		{Difficulty(99), "medium"}, // default for unknown
	}

	for _, tt := range tests {
		result := tt.difficulty.String()
		if result != tt.expected {
			t.Errorf("Difficulty(%d).String() = %q, expected %q", tt.difficulty, result, tt.expected)
		}
	}
}
