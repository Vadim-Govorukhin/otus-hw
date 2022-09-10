package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var timeout time.Duration

func main() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	arguments := os.Args
	argslen := len(arguments)
	if (argslen < 3) || (argslen > 4) {
		log.Printf("Usage of my TELNET client: [--timeout] host port\n got arguments %v", arguments[1:])
		os.Exit(1)
	}
	address := net.JoinHostPort(arguments[argslen-2], arguments[argslen-1])

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	err := client.Connect()
	if err != nil {
		log.Panicf("Connection error: %e", err)
	}
	defer client.Close()
	fmt.Fprintf(os.Stderr, "...Connected to %s with timeout %s\n", address, timeout)

	errorCh := make(chan error, 2)
	gracefulShoutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShoutdown, syscall.SIGINT, syscall.SIGQUIT, os.Interrupt)

	go func() {
		errorCh <- client.Send()
		log.Println("[sender] done")
	}()
	go func() {
		errorCh <- client.Receive()
		log.Println("[receiver] done")
	}()

	for {
		select {
		case err := <-errorCh:
			if err != nil {
				log.Panicf("got error %#v", err)
			}
			fmt.Fprintf(os.Stderr, "...EOF\n")
			return
		case <-gracefulShoutdown:
			return
		}
	}
}
