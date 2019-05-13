package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"./data"
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
	store := &FakeStore{}

	srv := createWriteServer(store)
	ts := httptest.NewServer(srv)
	defer ts.Close()

	http.Post(ts.URL, "text/plain", strings.NewReader("hello world 123"))

	if !reflect.DeepEqual(store.written, []string{"hello", "world", "123"}) {
		t.Error("Unexpected words written", store.written)
	}
}

func TestReadServer_GetOnly(t *testing.T) {
	srv := createReadServer(nil)
	ts := httptest.NewServer(srv)
	defer ts.Close()

	res, err := http.Post(ts.URL+"/stats", "text/plain", strings.NewReader("hello"))

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Error("Non-GET requests should be rejected")
	}
}

func TestReadServer_GetStats(t *testing.T) {
	store := FakeStore{}

	srv := createReadServer(&store)
	ts := httptest.NewServer(srv)
	defer ts.Close()

	res, _ := http.Get(ts.URL + "/stats")

	bodyBytes, _ := ioutil.ReadAll(res.Body)
	var actual data.Summary
	json.Unmarshal(bodyBytes, &actual)

	if !reflect.DeepEqual(actual, fakeSummary) {
		t.Error("Unexpected data returned", actual)
	}
}

type FakeStore struct {
	written []string
}

func (store *FakeStore) Write(word string) {
	store.written = append(store.written, word)
}

func (store *FakeStore) Query() data.Summary {
	return fakeSummary
}

var fakeSummary = data.Summary{
	Count:      123,
	TopWords:   []string{"damian"},
	TopLetters: []string{"d"},
}
