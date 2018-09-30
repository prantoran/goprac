package main

import (
  "fmt"
  "log"
  "net/http"
)

type pounds float32

func (p pounds) String() string {
  return fmt.Sprintf("£%.2f", p)
}

type database map[string]pounds

func (d database) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  switch r.URL.Path {
  case "/foo":
    fmt.Fprintf(w, "foo: %s\n", d["foo"])
  case "/bar":
    fmt.Fprintf(w, "bar: %s\n", d["bar"])
  default:
    w.WriteHeader(http.StatusNotFound)
    fmt.Fprintf(w, "No page found for: %s\n", r.URL)
  }
}

func main() {
  db := database{
    "foo": 1,
    "bar": 2,
  }

  log.Fatal(http.ListenAndServe("localhost:8000", db))
}
