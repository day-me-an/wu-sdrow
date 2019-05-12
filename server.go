package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"

	"./summary"
)

const (
	writeServerPort = 5555
	readServerPort  = 8080
)

func main() {
	agg := summary.New()

	// Must start in a separate goroutine because ListenAndServe blocks.
	go func() {
		fmt.Println("Starting read server on port", readServerPort)
		if err := http.ListenAndServe(fmt.Sprint(":", readServerPort), createReadServer(&agg)); err != nil {
			panic(err)
		}
	}()

	fmt.Println("Starting write server on port", writeServerPort)
	if err := http.ListenAndServe(fmt.Sprint(":", writeServerPort), createWriteServer(&agg)); err != nil {
		panic(err)
	}
}

// A server that listens for write requests.
func createWriteServer(agg summary.Aggregator) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// Efficiently iterate over the words as they come in.
		scanner := bufio.NewScanner(r.Body)
		scanner.Split(bufio.ScanWords)
		for scanner.Scan() {
			word := scanner.Text()
			agg.Write(word)
		}
	})

	return mux
}

// A server that provides stats.
func createReadServer(agg summary.Aggregator) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		summary := agg.Read()
		data, err := json.Marshal(summary)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})

	return mux
}