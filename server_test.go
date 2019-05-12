package main

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"./summary"
)

func TestWriteServer_PostOnly(t *testing.T) {
	srv := createWriteServer(nil)
	ts := httptest.NewServer(srv)
	defer ts.Close()

	res, err := http.Get(ts.URL)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Error("Non-POST requests should be rejected")
	}
}

func TestWriteServer_SubmitText(t *testing.T) {
	agg := FakeAggregator{}

	srv := createWriteServer(&agg)
	ts := httptest.NewServer(srv)
	defer ts.Close()

	http.Post(ts.URL, "text/plain", strings.NewReader("hello world 123"))

	if !reflect.DeepEqual(agg.written, []string{"hello", "world", "123"}) {
		t.Error("Unexpected words written", agg.written)
	}
}

type FakeAggregator struct {
	written []string
}

func (agg *FakeAggregator) Write(word string) {
	agg.written = append(agg.written, word)
}

func (agg *FakeAggregator) Read() summary.Summary {
	return summary.Summary{}
}
