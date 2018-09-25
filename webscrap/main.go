package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func ScrapeHTML() {
	resp, err := http.Get("https://github.com/trending")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
	}
	io.Copy(os.Stdout, resp.Body)
}

func main() {
	ScrapeHTML()
}
