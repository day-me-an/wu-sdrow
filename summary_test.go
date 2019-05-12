package summary

import (
	"reflect"
	"testing"
)

func TestWrite(t *testing.T) {
	agg := New()

	agg.Write("hello")
	agg.Write("hello")
	agg.Write("world")

	state := agg.Read()

	if state.Count != 2 {
		t.Error("Expected unique 2 words, but got", state.Count)
	}
}

func BenchmarkWrite(b *testing.B) {
	agg := New()

	// Don't include any time taken for initialisation in the benchmark.
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		// TODO: use random words here for a more realistic benchmark.
		agg.Write("hello")
	}
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
