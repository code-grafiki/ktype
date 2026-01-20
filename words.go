package main

import (
	"math/rand"
)

var wordList = []string{
	"the", "be", "to", "of", "and", "a", "in", "that", "have", "I",
	"it", "for", "not", "on", "with", "he", "as", "you", "do", "at",
	"this", "but", "his", "by", "from", "they", "we", "say", "her", "she",
	"or", "an", "will", "my", "one", "all", "would", "there", "their", "what",
	"so", "up", "out", "if", "about", "who", "get", "which", "go", "me",
	"when", "make", "can", "like", "time", "no", "just", "him", "know", "take",
	"people", "into", "year", "your", "good", "some", "could", "them", "see", "other",
	"than", "then", "now", "look", "only", "come", "its", "over", "think", "also",
	"back", "after", "use", "two", "how", "our", "work", "first", "well", "way",
	"even", "new", "want", "because", "any", "these", "give", "day", "most", "us",
	"is", "are", "was", "were", "been", "being", "has", "had", "does", "did",
	"very", "just", "should", "now", "here", "where", "why", "how", "when", "while",
	"each", "few", "more", "some", "such", "no", "nor", "too", "very", "can",
	"will", "just", "should", "now", "than", "then", "once", "here", "there", "when",
	"where", "why", "how", "all", "each", "every", "both", "few", "more", "most",
	"other", "some", "such", "only", "own", "same", "than", "too", "very", "just",
	"about", "after", "before", "between", "under", "again", "further", "once", "during", "through",
	"above", "below", "against", "without", "within", "along", "among", "around", "behind", "beyond",
	"never", "always", "often", "still", "already", "almost", "away", "back", "down", "home",
	"inside", "outside", "together", "apart", "instead", "rather", "simply", "quite", "enough", "much",
}

// getRandomWords returns n random words from the word list
func getRandomWords(n int) []string {
	words := make([]string, n)
	for i := 0; i < n; i++ {
		words[i] = wordList[rand.Intn(len(wordList))]
	}
	return words
}
