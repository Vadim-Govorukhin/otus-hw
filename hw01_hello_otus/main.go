package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	s := "Hello, OTUS!"
	fmt.Println(stringutil.Reverse(s))
}
