package main

import (
	"bufio"
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
	err = c.conn.Close()
	fmt.Fprintln(os.Stderr, "...Connection was closed by peer")
	return
}

func (c *Client) Receive() (err error) {
	for {
		data, err := bufio.NewReader(c.conn).ReadString('\n')
		log.Printf("[cl] receive '%s'\n", data)

		if err != nil {
			break
		}
		_, err = c.out.Write([]byte(data))
		if err != nil {
			break
		}
	}
	return
}

func (c *Client) Send() (err error) {
	for {
		data, err := bufio.NewReader(c.in).ReadString('\n')
		log.Printf("[cl] read %#v with error %e\n", data, err)
		log.Printf("[cl] compare %#v with %#v\n", data, `\x04\r\n`)
		if data == `\x04\r\n` {
			err = io.EOF
		}
		if err != nil {
			log.Printf("[cl] return with error %e\n", err)
			break
		}
		_, err = fmt.Fprint(c.conn, data)
		if err != nil {
			break
		}
		log.Printf("[cl] send %#v with error %e\n", data, err)
	}
	return
}

// Place your code here.
// P.S. Author's solution takes no more than 50 lines.
