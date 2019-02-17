package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/protected", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, I'm protected"))
	})

	log.Fatal(http.ListenAndServe(":9096", nil))
}
