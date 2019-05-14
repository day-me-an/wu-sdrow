package data

import (
	"strings"
	"sync"
	"unicode"

	"github.com/wangjia184/sortedset"
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

	words   *sortedset.SortedSet
	letters *sortedset.SortedSet
}

func NewMutexStore() *MutexStore {
	return &MutexStore{
		words:   sortedset.New(),
		letters: sortedset.New(),
	}
}

func (store *MutexStore) Write(word string) {
	// Words should not be discriminated by case.
	word = strings.ToLower(word)

	store.lock.Lock()
	defer store.lock.Unlock()

	store.words.AddOrUpdate(word, 1, nil)

	for _, char := range word {
		if unicode.IsLetter(char) {
			letter := string(char)

			store.letters.AddOrUpdate(letter, 1, nil)
		}
	}
}

func (store *MutexStore) Query() Summary {
	store.lock.RLock()
	defer store.lock.RUnlock()

	return Summary{
		Count:      store.words.GetCount(),
		TopWords:   sortedSetTopN(store.words, 5),
		TopLetters: sortedSetTopN(store.letters, 5),
	}
}

func sortedSetTopN(ss *sortedset.SortedSet, n int) []string {
	nodes := ss.GetByRankRange(1, 5, false)

	var items []string

	for _, node := range nodes {
		items = append(items, node.Key())
	}

	return items
}
