package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/hello", Hello)
	http.HandleFunc("/healthz", Healthz)

	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		panic(err)
	}
}

func Healthz(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("It's alive!"))
	}
}

func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, Master!")
}
