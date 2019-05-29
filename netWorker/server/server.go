package server

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
)

type Server struct {
	Conn chan net.Conn
}

func (s *Server) Run(addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		fmt.Println("有一个客户端上线：", conn.RemoteAddr().String())

		s.Conn <- conn

	}

}

func (s *Server) Read() bytes.Buffer {
	fmt.Println("server")

	conn := <-s.Conn

	defer conn.Close()
	buf, err := ioutil.ReadAll(conn)
	if err != nil {
		panic(err)
	}
	fmt.Println(buf)
	var bb bytes.Buffer
	_, err = bb.Write(buf)
	if err != nil {
		panic(err)
	}
	return bb

}

func (s *Server) Write(bytes.Buffer) {}
