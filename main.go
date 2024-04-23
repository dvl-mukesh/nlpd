package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/dvl-mukesh/nlp"
)

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/tokenize", tokenizeHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}

func tokenizeHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	lr := io.LimitReader(r.Body, 1_000_000)
	bytes, err := io.ReadAll(lr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	text := string(bytes)

	if len(text) == 0 {
		http.Error(w, "empty text", http.StatusBadRequest)
		return
	}
	tokens := nlp.Tokenize(text)
	response, err := json.Marshal(tokens)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(response))
}
