package main

import (
	"bytes"
	"flag"
	"fmt"
	"net"

	"test/fileTest2/gofile/internal"
	"test/fileTest2/gofile/networker/client"
	"test/fileTest2/gofile/networker/server"
)

var f = flag.String("f", "test.txt", "input the file name")
var hostname = flag.String("s", "", "input server ip")
var port = flag.String("p", "", "input server port")

type Control interface {
	Reader() bytes.Buffer
	Writer(bb bytes.Buffer)
}

type Control2 interface {
	Reader(*net.Conn) bytes.Buffer
}

func init() {
	flag.Parse()
}

func main() {
	var ctrl Control
	if *hostname != "" {
		fmt.Println("send file")
		c := client.New(*hostname)
		bbytes := internal.Defaultbuffer.GetBytesbuffer(*f)

		ctrl = c
		ctrl.Writer(bbytes)
	} else if *port != "" {
		fmt.Println("get file")
		var s = server.Server{}
		go s.Run(*port)
		ctrl = s
		for {
			bb := ctrl.Reader()
			internal.Defaultbuffer.PutBytesbufferToFile(bb.Bytes())
		}

	}

}
