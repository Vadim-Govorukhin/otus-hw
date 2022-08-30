package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var timeout time.Duration

func main() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	arguments := os.Args
	if len(arguments) != 4 {
		log.Printf("Usage: %s [--timeout] host port", arguments[0])
		os.Exit(1)
	}
	address := net.JoinHostPort(arguments[len(arguments)-2], arguments[len(arguments)-1])

	t := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	err := t.Connect()
	if err != nil {
		log.Println("Connection error: ", err)
		return
	}

	fmt.Fprintf(os.Stderr, "...Connected to %s with timeout %s\n", address, timeout)
	gracefulShutdown := make(chan os.Signal, 1) //////
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM)

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
	log.Println("Close client")
	t.Close()

	// Place your code here,
	// P.S. Do not rush to throw context down, think think if it is useful with blocking operation?
}
