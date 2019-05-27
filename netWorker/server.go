package networker

import (
	"fmt"
	"io/ioutil"
	"net"
	"test/fileTest2/gofile/internal"
)

type Server struct{}

func (s Server) Run(addr string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		fmt.Println("有一个客户端上线：", conn.RemoteAddr().String())
		go s.Reader(conn)
	}
}

func (s Server) Reader(conn net.Conn) {
	defer conn.Close()
	buf, err := ioutil.ReadAll(conn)
	if err != nil {
		panic(err)
	}
	internal.Defaultbuffer.PutBytesbufferToFile(buf)
}
