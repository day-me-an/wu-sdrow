package main

import (
	"fmt"
	"net/http"
)

func main() {
	if err := http.ListenAndServe(fmt.Sprint(":", 5555), createWriteServer()); err != nil {
		panic(err)
	}
}

// A server that listens for write requests.
func createWriteServer() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		w.Write([]byte("hello world"))
	})

	return mux
}
