package main

import (
	"log"
	"reflect"
)

type Gift struct {
	Sender    string
	Recipient string
	Number    uint
	Contents  string
}

func main() {
	g := Gift{
		Sender:    "Hank",
		Recipient: "Sue",
		Number:    1,
		Contents:  "Scarf",
	}

	t := reflect.TypeOf(g)

	if kind := t.Kind(); kind != reflect.Struct {
		log.Fatalf("This program expects to work on a struct; we got a %v instead.", kind)
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		log.Printf("Field %03d: %-10.20s %-10.20v size: %4v", i, f.Name, f.Type.Kind(), f.Type.Size())
	}
}

/*
2009/11/10 23:00:00 Field 000: Sender     string     size:    8
2009/11/10 23:00:00 Field 001: Recipient  string     size:    8
2009/11/10 23:00:00 Field 002: Number     uint       size:    4
2009/11/10 23:00:00 Field 003: Contents   string     size:    8
*/
