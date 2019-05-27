package main

import (
	"bytes"
	"flag"
	"net"

	// "test/fileTest2/gofile/internal"
	"test/fileTest2/gofile/networker"
)

var f = flag.String("f", "test.txt", "input the file name")
var hostname = flag.String("s", "127.0.0.1:5000", "input server ip")
var port = flag.String("p", ":5000", "input server port")

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
	// w := networker.New(*hostname)
	// bbytes := internal.Defaultbuffer.GetBytesbuffer(*f)
	// var ctrl Control
	// ctrl = w
	// ctrl.Writer(bbytes)
	var s = networker.Server{}
	s.Run(*port)
}
