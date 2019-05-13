package data

import (
	"reflect"
	"testing"
)

func TestWrite(t *testing.T) {
	store := NewMutexStore()
	store.Write("hello")
	// No checks.
}

func BenchmarkWrite(b *testing.B) {
	store := NewMutexStore()

	// Don't include any time taken for initialisation in the benchmark.
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		// TODO: use random words here for a more realistic benchmark.
		store.Write("hello")
	}
}

func TestQuery_Words(t *testing.T) {
	summary := summaryFrom("hello", "hello", "world")

	if summary.Count != 2 {
		t.Error("Expected unique 2 words, but got", summary.Count)
	}

	if !reflect.DeepEqual(summary.TopWords, []string{"hello", "world"}) {
		t.Error("Wrong top words", summary.TopWords, len(summary.TopWords))
	}
}

func TestQuery_Letters(t *testing.T) {
	summary := summaryFrom("aaab")

	if !reflect.DeepEqual(summary.TopLetters, []string{"a", "b"}) {
		t.Error("Wrong top letters", summary.TopLetters, len(summary.TopLetters))
	}
}

func TestQuery_Casing(t *testing.T) {
	summary := summaryFrom("damian", "Damian", "DAMIAN")

	if !reflect.DeepEqual(summary.TopWords, []string{"damian"}) {
		t.Error("Wrong top words", summary.TopWords, len(summary.TopWords))
	}
}

// Helper function that creates a data store, writes the words and performs a query.
func summaryFrom(words ...string) Summary {
	store := NewMutexStore()

	for _, word := range words {
		store.Write(word)
	}

	return store.Query()
}

func TestTopN(t *testing.T) {
	m := map[string]int{"bob": -4, "hello": 5, "world": 8}

	actual := topN(m, 2)
	expected := []string{"world", "hello"}

	if !reflect.DeepEqual(actual, expected) {
		t.Error("Top items not as expected", actual)
	}
}

func TestTopMismatchedSize(t *testing.T) {
	m := map[string]int{"hello": 5, "world": 8}

	actual := topN(m, 3)
	expected := []string{"world", "hello"}

	if len(actual) != len(expected) {
		t.Error("Top items length should be less than or equal to the actual map length")
	}
}
