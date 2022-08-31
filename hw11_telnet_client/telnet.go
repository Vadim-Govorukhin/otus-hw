package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type Client struct {
	conn    net.Conn
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Client{nil, address, timeout, in, out}
}

func (c *Client) Connect() (err error) {
	c.conn, err = net.DialTimeout("tcp", c.address, c.timeout)
	return err
}

func (c *Client) Close() (err error) {
	err = c.in.Close()
	if err != nil {
		log.Printf("error while closing 'in' %e", err)
	}

	err = c.conn.Close()
	if err != nil {
		log.Printf("error while closing connection %e", err)
	}
	fmt.Fprintln(os.Stderr, "...Connection was closed by peer")
	return
}

func (c *Client) Receive() error {
	_, err := io.Copy(c.out, c.conn)
	return err
}

func (c *Client) Send() error {
	_, err := io.Copy(c.conn, c.in)
	return err
}

// Place your code here.
// P.S. Author's solution takes no more than 50 lines.
