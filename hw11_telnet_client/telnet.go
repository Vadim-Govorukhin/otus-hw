package main

import (
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type Telnet struct {
	conn    net.Conn
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return Telnet{nil, address, timeout, in, out}
}

func (t Telnet) Connect() (err error) {
	t.conn, err = net.DialTimeout("tcp", t.address, t.timeout)
	return err
}

func (t Telnet) Close() (err error) {
	err = t.conn.Close()
	return
}

func (t Telnet) Receive() (err error) {
	//_, err = t.in.Read("\n")
	return
}

func (t Telnet) Send() (err error) {
	_, err = t.out.Write([]byte("Received"))
	return
}

// Place your code here.
// P.S. Author's solution takes no more than 50 lines.
