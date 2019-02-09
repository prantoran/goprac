package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handle)

	log.Print("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func handle(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, Guest"))
}
