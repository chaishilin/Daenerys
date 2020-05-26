package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	fmt.Println(isTitle("dsfghf:dgfh"))
	fmt.Println(isTitle("dsffh"))
}
func isTitle(str string) bool {
	for len(str) > 0 {
		r, size := utf8.DecodeRuneInString(str)
		if r == ':'{
			return false
		}
		str = str[size:]
	}
	return true
}