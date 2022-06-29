package main

import "fmt"

func main() {
	s := "Привет"

	for _, c := range s {
		fmt.Printf("%c\n", c)
	}
}
