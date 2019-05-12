package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteServer_PostOnly(t *testing.T) {
	srv := createWriteServer()
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
