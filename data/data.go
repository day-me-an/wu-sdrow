package data

import (
	"strings"
	"sync"
	"unicode"
)

type Summary struct {
	// The number of _unique_ words seen.
	Count int `json:"count"`
	// The top N occurring words seen.
	TopWords []string `json:"top_words"`
	// The top N letters seen in all words.
	TopLetters []string `json:"top_letters"`
}

type Store interface {
	Write(string)
	Query() Summary
}

type MutexStore struct {
	// Chosen a RWMutex over a regular Mutex because multiple simultaneous reads don't need to lock unless sonething is being written.
	lock sync.RWMutex

	words   map[string]int
	letters map[string]int
}

func NewMutexStore() *MutexStore {
	return &MutexStore{
		words:   make(map[string]int),
		letters: make(map[string]int),
	}
}

func (store *MutexStore) Write(word string) {
	// Words should not be discriminated by case.
	word = strings.ToLower(word)

	store.lock.Lock()
	defer store.lock.Unlock()

	if count, exists := store.words[word]; exists {
		store.words[word] = count + 1
	} else {
		store.words[word] = 1
	}

	for _, char := range word {
		if unicode.IsLetter(char) {
			letter := string(char)

			if count, exists := store.letters[letter]; exists {
				store.letters[letter] = count + 1
			} else {
				store.letters[letter] = 1
			}
		}
	}
}

func (store *MutexStore) Query() Summary {
	store.lock.RLock()
	defer store.lock.RUnlock()

	return Summary{
		Count:      len(store.words),
		TopWords:   topN(store.words, 5),
		TopLetters: topN(store.letters, 5),
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
