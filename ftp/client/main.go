package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jlaffaye/ftp"
)

func connect() (client *ftp.ServerConn, err error) {
	client, err = ftp.Dial("localhost:2001")
	if err != nil {
		return nil, err
	}

	if err := client.Login("yo", "yo"); err != nil {
		return nil, err
	}

	return client, nil
}

func filenames(root string) []string {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".json" {
			return nil
		}

		files = append(files, path)
		// files = append(files, info.Name())

		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fmt.Println(file)
	}

	return files
}

type filedata struct {
	name string
	data []byte
}

func jsonFiles() ([]filedata, error) {

	filenames := filenames("test")

	files := []filedata{}

	for _, u := range filenames {
		dat, err := ioutil.ReadFile(u)

		if err != nil {
			return nil, err
		}

		v := strings.Split(u, "/")

		files = append(files, filedata{
			name: v[len(v)-1],
			data: dat,
		})

	}

	return files, nil

}

func main() {

	client, err := connect()
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(client)

	pwd, err := client.CurrentDir()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("pwd:", pwd)

	err = client.MakeDir("ftpclientmkdir")
	if err != nil {
		log.Fatal(err)
	}

	// err = client.RemoveDir("ftpclientmkdir2")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	jsonFiles, err := jsonFiles()

	for _, u := range jsonFiles {

		fmt.Println("u.name:", u.name, " data:", string(u.data))

		r := bytes.NewReader(u.data)

		err := client.Stor("ftpclientmkdir/"+u.name, r)
		if err != nil {
			log.Fatal(err)
		}
	}

	entries, _ := client.List("*")

	for _, entry := range entries {
		name := entry.Name
		// reader, err := client.Retr(name)
		// if err != nil {
		// 	panic(err)
		// }
		// client.Delete(name)
		fmt.Println("name:", name)
	}
}
