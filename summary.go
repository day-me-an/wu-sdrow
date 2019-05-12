package summary

import (
	"sync"
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

	words map[string]int
}

func New() Aggregator {
	return Aggregator{
		words: make(map[string]int),
	}
}

func (agg *Aggregator) Write(word string) {
	agg.lock.Lock()
	defer agg.lock.Unlock()

	if count, exists := agg.words[word]; exists {
		agg.words[word] = count + 1
	} else {
		agg.words[word] = 1
	}
}

func (agg *Aggregator) Read() Summary {
	agg.lock.RLock()
	defer agg.lock.RUnlock()

	return Summary{
		Count: len(agg.words),
	}
}
