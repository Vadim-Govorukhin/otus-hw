package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	var timeout string

	flag.StringVar(&timeout, "timeout", "10s", "timeout of connection")
	flag.Parse()

	dur, err := time.ParseDuration(timeout)
	if err != nil {
		log.Fatalf("Cannot parse timeout: %v", err)
	}

	arguments := os.Args
	fmt.Println(arguments)
	if len(arguments) != 3 {
		log.Printf("Usage: %s host port ", os.Args[0])
		os.Exit(1)
	}

	HOST := arguments[1]
	PORT := arguments[2]
	address := HOST + ":" + PORT
	t := NewTelnetClient(address, dur, os.Stdin, os.Stdout)
	t.Connect()

	// Place your code here,
	// P.S. Do not rush to throw context down, think think if it is useful with blocking operation?
}
