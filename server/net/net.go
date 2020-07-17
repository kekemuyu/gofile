package net

import (
	"bytes"
	"net"

	"github.com/kekemuyu/gofile/msg"

	"github.com/kekemuyu/gofile/handler"

	"github.com/donnie4w/go-logger/logger"
)

type Server struct {
}

var DefaultServer = Server{}

func (s *Server) Run(addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Error(err)
			continue
		}
		logger.Debug("有一个客户端上线：", conn.RemoteAddr().String())
		hlr := handler.Handler{
			Rwc: conn,
		}
		go hlr.HandleLoop() //客户端连接处理，读写数据
	}

}

func (s *Server) Read(b []byte) (n int, err error) {
	return
}

func (s *Server) Write(b []byte) (n int, err error) {
	return
}

func (s *Server) WriteHandle(bytes.Buffer) {

}

func (s *Server) ReadHandle(msg msg.Msg) {
	// switch msg.Id {
	// case 0x00000001:
	// 	filename := string(msg.Data)
	// 	var err error
	// 	internal.Default, err = internal.Newfile(filename)
	// 	if err != nil {
	// 		logger.Error("newfile err:", err)
	// 		return
	// 	}
	// 	logger.Debug("newfile success:", internal.Default, err)
	// case 0x00000002:
	// 	internal.Default.Write(msg.Data)
	// 	if int64(len(msg.Data)) == internal.Default.Blocksize {
	// 		internal.Default.Writeoffset += internal.Default.Blocksize
	// 	} else {
	// 		internal.Default.Writeoffset += int64(len(msg.Data))
	// 	}
	// case 0x00000003:
	// 	internal.Default.Write(msg.Data)
	// case 0x00000004:

	// default:
	// }
}
