package main

import (
	"io"
	"log"
	"net"
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
	log.Printf("...Connected to %s\n", c.address)
	return err
}

func (c *Client) Close() (err error) {
	err = c.conn.Close()
	log.Println("...Connection was closed by")
	return
}

func (c *Client) Receive() (err error) {
	//_, err = c.in.Read("\n")
	return
}

func (c *Client) Send() (err error) {
	_, err = c.out.Write([]byte("Received"))
	return
}

// Place your code here.
// P.S. Author's solution takes no more than 50 lines.
