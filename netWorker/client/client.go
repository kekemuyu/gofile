package client

import (
	"bytes"
	"fmt"
	"net"
)

type Client struct {
	Conn net.Conn
}

func New(hostname string) *Client {
	conn, err := net.Dial("tcp", hostname)
	if err != nil {
		panic(conn)
	}

	return &Client{
		Conn: conn,
	}
}

func (c *Client) Writer(bb bytes.Buffer) {
	defer c.Conn.Close()
	_, err := c.Conn.Write(bb.Bytes())
	if err != nil {
		fmt.Println(err)
	}
}

func (c *Client) Reader() bytes.Buffer {
	return bytes.Buffer{}
}
