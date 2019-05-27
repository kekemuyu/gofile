package networker

import (
	"bytes"
	"fmt"
	"net"
)

type Worker struct {
	Conn *net.TCPConn
}

func New(hostName string) *Worker {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", hostName)
	if err != nil {
		panic(err)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	defer conn.Close()
	if err != nil {
		panic(conn)
	}

	return &Worker{
		Conn: conn,
	}
}

func (w *Worker) Writer(bb bytes.Buffer) {
	_, err := w.Conn.Write(bb.Bytes())
	if err != nil {
		fmt.Println(err)
	}
}

func (w *Worker) Reader() bytes.Buffer {

	w.Conn.re
}
