package main

import (
	"bufio"
	"fmt"
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
	return err
}

func (c *Client) Close() (err error) {
	err = c.conn.Close()
	log.Println("...Connection was closed by")
	return
}

func (c *Client) Receive() (err error) {
	for {
		data, err := bufio.NewReader(c.conn).ReadString('\n')
		log.Println("[cl] receive ", data)
		if err != nil {
			break
		}
		c.out.Write([]byte(data))
	}
	return
}

func (c *Client) Send() (err error) {
	for {
		data, err := bufio.NewReader(c.in).ReadString('\n')
		fmt.Fprint(c.conn, data)
		log.Println("[cl] send ", data)
		if err != nil {
			break
		}
	}
	return
}

// Place your code here.
// P.S. Author's solution takes no more than 50 lines.
