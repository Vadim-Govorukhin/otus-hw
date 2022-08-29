package main

import (
	"flag"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "timeout of connection")
	flag.Parse()

	arguments := os.Args
	if len(arguments) != 3 {
		log.Printf("Usage: %s host port ", os.Args[0])
		os.Exit(1)
	}

	HOST := arguments[1]
	PORT := arguments[2]
	address := net.JoinHostPort(HOST, PORT)

	t := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	t.Connect()

	// Place your code here,
	// P.S. Do not rush to throw context down, think think if it is useful with blocking operation?
}
