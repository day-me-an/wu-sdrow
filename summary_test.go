package summary

import (
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
