package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("hello world")
	for _, i := range os.Args[1:] {
		fmt.Println(i)
	}
}

func hello() int {
	fmt.Println("hello")
	return 1
}
