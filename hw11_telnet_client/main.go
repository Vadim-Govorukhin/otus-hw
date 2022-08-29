package main

import (
	"flag"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

func main() {
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "timeout of connection")
	flag.Parse()

	arguments := os.Args
	if len(arguments) != 3 {
		log.Printf("Usage: %s host port", arguments[0])
		os.Exit(1)
	}
	address := net.JoinHostPort(arguments[1], arguments[2])

	t := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	err := t.Connect()
	if err != nil {
		log.Println("Connection error: ", err)
		return
	}
	log.Printf("...Connected to %s\n", address)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		err := t.Send()
		if err != nil {
			log.Println("[sender] error: ", err)
			return
		}
	}()
	go func() {
		defer wg.Done()
		err := t.Receive()
		if err != nil {
			log.Println("[receiver] error: ", err)
			return
		}
	}()

	wg.Wait()
	t.Close()

	// Place your code here,
	// P.S. Do not rush to throw context down, think think if it is useful with blocking operation?
}
