package main

import (
	"fmt"
	"log"

	"go.uber.org/dig"
)

type car struct {
	size   int
	wheels int
}

type ford struct {
	size   int
	wheels int
	speed  int
}

func main() {
	c := dig.New()

	if err := c.Provide(func() *car {
		return &car{
			size:   10,
			wheels: 4,
		}
	}); err != nil {
		log.Println(err)
	}

	if err := c.Provide(func(c *car) *ford {
		return &ford{
			size:   c.size,
			wheels: c.wheels,
			speed:  20,
		}
	}); err != nil {
		log.Println(err)
	}

	if err := c.Invoke(func(f *ford) error {
		fmt.Printf("ford: %#v\n", f)
		return nil
	}); err != nil {
		log.Println(err)
	}
}
