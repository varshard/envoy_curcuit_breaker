package main

import (
	"bytes"
	"fmt"

	//"math/rand"
	"net/http"
)

func main() {
	tryCount := 0

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "pong")
	})

	http.HandleFunc("/failed", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "failed", 500)
	})

	http.HandleFunc("/retry", func(w http.ResponseWriter, request *http.Request) {
		tryCount++
		if tryCount%3 != 0 {
			http.Error(w, fmt.Sprintf("retry this %d", tryCount), 500)
			return
		}

		fmt.Fprintf(w, "Finally %d", tryCount)
	})

	http.HandleFunc("/google", func(w http.ResponseWriter, r *http.Request) {
		_, err := http.NewRequest("get", "envoy:10000", new(bytes.Buffer))
		if err != nil {
			http.Error(w, "failed", 500)
		} else {
			fmt.Fprint(w, "googled")
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "up")
	})

	http.ListenAndServe(":6060", nil)
}
