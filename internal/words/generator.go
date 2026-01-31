package words

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GetRandom returns n random words from the word list for given difficulty
// Ensures no word appears consecutively for better typing flow
func GetRandom(n int, difficulty Difficulty) []string {
	if n <= 0 {
		return []string{}
	}

	wordList := GetList(difficulty)
	words := make([]string, n)
	for i := 0; i < n; i++ {
		word := wordList[rand.Intn(len(wordList))]
		// Prevent consecutive duplicates
		if i > 0 && word == words[i-1] && len(wordList) > 1 {
			// Pick a different word
			for word == words[i-1] {
				word = wordList[rand.Intn(len(wordList))]
			}
		}
		words[i] = word
	}
	return words
}

// GetList returns the word list for a given difficulty
func GetList(d Difficulty) []string {
	switch d {
	case DifficultyEasy:
		return easyWords
	case DifficultyHard:
		return hardWords
	default:
		return mediumWords
	}
}

// GetRandomWithComplexity returns words with optional punctuation/numbers
func GetRandomWithComplexity(n int, difficulty Difficulty, complexity Complexity) []string {
	if n <= 0 {
		return []string{}
	}

	baseWords := GetRandom(n, difficulty)
	if complexity == ComplexityNormal {
		return baseWords
	}

	words := make([]string, n)

	for i := 0; i < n; i++ {
		word := baseWords[i]

		switch complexity {
		case ComplexityPunctuation:
			// 30% chance to add punctuation to a word
			if rand.Float32() < 0.3 {
				word = AddPunctuation(word)
			}
		case ComplexityNumbers:
			// 20% chance to replace word with number, 10% to add number to word
			r := rand.Float32()
			if r < 0.2 {
				word = numberList[rand.Intn(len(numberList))]
			} else if r < 0.3 {
				word = word + numberList[rand.Intn(len(numberList))]
			}
		case ComplexityFull:
			// Mix of punctuation and numbers
			r := rand.Float32()
			if r < 0.25 {
				// Replace with number
				word = numberList[rand.Intn(len(numberList))]
			} else if r < 0.45 {
				// Add punctuation
				word = AddPunctuation(word)
			} else if r < 0.55 {
				// Number-punctuation combo
				word = numberPunctuationCombos[rand.Intn(len(numberPunctuationCombos))]
			}
		}

		words[i] = word
	}

	return words
}

// AddPunctuation adds punctuation to a word
func AddPunctuation(word string) string {
	// Different punctuation strategies
	switch rand.Intn(5) {
	case 0:
		// Add trailing punctuation
		return word + punctuationMarks[rand.Intn(len(punctuationMarks))]
	case 1:
		// Add leading punctuation
		return punctuationMarks[rand.Intn(len(punctuationMarks))] + word
	case 2:
		// Wrap in punctuation
		p1 := punctuationMarks[rand.Intn(len(punctuationMarks))]
		p2 := punctuationMarks[rand.Intn(len(punctuationMarks))]
		return p1 + word + p2
	case 3:
		// Add apostrophe combo
		return word + punctuationCombos[rand.Intn(6)+8] // 's, n't, 're, 'll, 'd, 've, 'm
	default:
		// Add internal punctuation
		if len(word) > 2 {
			mid := len(word) / 2
			return word[:mid] + punctuationMarks[rand.Intn(len(punctuationMarks))] + word[mid:]
		}
		return word + punctuationMarks[rand.Intn(len(punctuationMarks))]
	}
}
