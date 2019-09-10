package net

import (
	"bytes"
	"gofile/internal"
	"gofile/msg"
	"net"

	"github.com/donnie4w/go-logger/logger"
)

type Server struct {
	Conn net.Conn
	Fman internal.Fileman
}

var DefaultServer = new(Server)

func (s *Server) Run(addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	conn, err := listener.Accept()
	if err != nil {
		panic(err)
	}
	logger.Debug("有一个客户端上线：", conn.RemoteAddr().String())

	DefaultServer.Conn = conn

}

func (s *Server) Read(b []byte) (n int, err error) {
	return DefaultServer.Conn.Read(b)
}

func (s *Server) Write(b []byte) (n int, err error) {
	return
}

func (s *Server) WriteHandle(bytes.Buffer) {

}

func (s *Server) ReadHandle(msg msg.Msg) {
	switch msg.Id {
	case 0x00000001:
		filename := string(msg.Data)
		var err error
		s.Fman, err = internal.Newfile(filename)
		if err != nil {
			logger.Error("newfile err:", err)
			return
		}
		logger.Debug("newfile success:", s.Fman, err)
	case 0x00000002:
		s.Fman.Write(msg.Data)
		s.Fman.Writeoffset += s.Fman.Blocksize
	case 0x00000003:
		s.Fman.Write(msg.Data)
	default:
	}
}
