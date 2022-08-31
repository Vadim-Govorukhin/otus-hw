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

const argslen = 4

var timeout time.Duration

func main() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	arguments := os.Args
	if len(arguments) != argslen {
		log.Printf("Usage: %s [--timeout] host port", arguments[0])
		os.Exit(1)
	}
	address := net.JoinHostPort(arguments[argslen-2], arguments[argslen-1])

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	err := client.Connect()
	if err != nil {
		log.Println("Connection error: ", err)
		return
	}

	fmt.Fprintf(os.Stderr, "...Connected to %s with timeout %s\n", address, timeout)
	//ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGQUIT, os.Interrupt)
	errorCh := make(chan error, 2)
	gracefulShoutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShoutdown, syscall.SIGINT, syscall.SIGQUIT, os.Interrupt)

	defer func() {
		log.Println("Close client")
		client.Close()
	}()

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
				log.Printf("got error %#v", err)
				return
			}

			return
		case <-gracefulShoutdown:
			return
		}
	}

	// Place your code here,
	// P.S. Do not rush to throw context down, think think if it is useful with blocking operation?
}
