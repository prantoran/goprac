package main

import (
  "fmt"
  "os"

  "path/filepath"

)

func main() {
  dir, err := os.Open("../")
  if err != nil {
    return
  }
  defer dir.Close()

  fileInfos, err := dir.Readdir(-1)
  if err != nil {
    return
  }
  for _, fi := range fileInfos {
    fmt.Println(fi.Name())
  }
	fmt.Println("******")
  filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
    fmt.Println(path)
    return nil
  })
}
