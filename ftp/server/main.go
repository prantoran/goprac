package main

import (
	"fmt"

	filedriver "github.com/goftp/file-driver"
	"github.com/goftp/server"
)

func main() {

	// cmd := exec.Command("pwd")
	// var out bytes.Buffer
	// cmd.Stdout = &out
	// err := cmd.Run()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// filepath := out.String()

	// if filepath[len(filepath)-1] == '\n' {
	// 	filepath = filepath[0 : len(filepath)-1]
	// }

	filepath := "/Users/pinku/work/src/github.com/prantoran/ftp"

	fmt.Printf("filepath: %q\n", filepath)

	factory := &filedriver.FileDriverFactory{
		RootPath: filepath,
		Perm:     server.NewSimplePerm("yo", "yo"),
	}

	// Perm field defines how the user are going to be authenticated

	opts := &server.ServerOpts{
		Factory:  factory,
		Port:     2001,
		Hostname: "127.0.0.1",
		Auth: &server.SimpleAuth{
			Name:     "yo",
			Password: "yo",
		},
	}
	server := server.NewServer(opts)
	server.ListenAndServe()

}
