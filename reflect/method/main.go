// Examining the itab/Method-Set
// access portions of the itab to extract method information from the interface

package main

import (
	"log"
	"reflect"
)

type Reindeer string

func (r Reindeer) TakeOff() {
	log.Printf("%q lifts off.", r)
}

func (r Reindeer) Land() {
	log.Printf("%q gently lands.", r)
}

func (r Reindeer) ToggleNose() {
	if r != "rudolph" {
		panic("invalid reindeer operation")
	}
	log.Printf("%q nose changes state.", r)
}

func main() {
	r := Reindeer("rudolph")

	t := reflect.TypeOf(r)

	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		log.Printf("%s", m.Name)
	}

	func() {
		x, ok := t.MethodByName("TakeOff")

		if ok == false {
			log.Println("method not found")
			return
		}

		xFunc := x.Func
		xIndx := x.Index

		log.Println("x.Func:", xFunc, "\txIndx:", xIndx)

	}()

	/*
		2019/01/03 16:00:33 Land
		2019/01/03 16:00:33 TakeOff
		2019/01/03 16:00:33 ToggleNose
		2019/01/03 16:00:33 x.Func: 0x10ad2c0   xIndx: 1
	*/
}
