package summary

import (
	"strings"
	"sync"
	"unicode"
)

type Summary struct {
	// The number of _unique_ words seen.
	Count int
	// The top N occurring words seen.
	TopWords []string
	// The top N letters seen in all words.
	TopLetters []string
}

type Aggregator struct {
	// Chosen a RWMutex over a regular Mutex because multiple simultaneous reads don't need to lock unless sonething is being written.
	lock sync.RWMutex

	words   map[string]int
	letters map[string]int
}

func New() Aggregator {
	return Aggregator{
		words:   make(map[string]int),
		letters: make(map[string]int),
	}
}

func (agg *Aggregator) Write(word string) {
	// Words should not be discriminated by case.
	word = strings.ToLower(word)

	agg.lock.Lock()
	defer agg.lock.Unlock()

	if count, exists := agg.words[word]; exists {
		agg.words[word] = count + 1
	} else {
		agg.words[word] = 1
	}

	for _, char := range word {
		if unicode.IsLetter(char) {
			letter := string(char)

			if count, exists := agg.letters[letter]; exists {
				agg.letters[letter] = count + 1
			} else {
				agg.letters[letter] = 1
			}
		}
	}
}

func (agg *Aggregator) Read() Summary {
	agg.lock.RLock()
	defer agg.lock.RUnlock()

	return Summary{
		Count:      len(agg.words),
		TopWords:   topN(agg.words, 5),
		TopLetters: topN(agg.letters, 5),
	}
}

// Finds the N highest value entries in a map.
func topN(m map[string]int, n int) []string {
	// The map size may be less than the requested N.
	if n > len(m) {
		n = len(m)
	}

	top := make([]string, n)
	picked := make(map[string]bool)

	for i := 0; i < n; i++ {
		highestSoFar := 0

		for item, count := range m {
			if picked[item] {
				continue
			}

			if count > highestSoFar {
				highestSoFar = count
				top[i] = item
			}
		}

		picked[top[i]] = true
	}

	return top
}
