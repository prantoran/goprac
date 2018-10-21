package main

import (
	"fmt"
	"regexp"
	"strings"
)

func parseURL(s string) []string {
	re := regexp.MustCompile(`\{(.*?)\}`)
	// fmt.Printf("Pattern: %v\n", re.String())      // print pattern
	// fmt.Println("Matched:", re.MatchString(s)) // true

	submatchall := re.FindAllString(s, -1)

	if submatchall == nil {
		fmt.Println("en")
		return nil
	}
	ret := []string{}
	for _, element := range submatchall {
		element = strings.Trim(element, "{")
		element = strings.Trim(element, "}")
		ret = append(ret, element)
	}

	return ret
}

func main() {
	// str1 := "this/is/a/{sample}/{{string}}/with/{SOME}/special/words"

	ret := parseURL("absc")

	if ret == nil {
		fmt.Println("en2")
	}

	fmt.Println("\nText between square brackets:", ret)

}
