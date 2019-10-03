package main

import (
	"bytes"
	"fmt"
	//"math/rand"
	"net/http"
)

func main() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		//rand.Seed(time.Now().UnixNano())
		//sleep := rand.Int31n(10)
		//time.Sleep(time.Duration(sleep) * time.Second)

		fmt.Fprint(w, "pong")
	})

	http.HandleFunc("/failed", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "failed", 500)
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
