package main

import (
	"flag"
	"fmt"
	"net"

	"test/fileTest2/gofile/internal"
	"test/fileTest2/gofile/networker/client"
	"test/fileTest2/gofile/networker/server"
	"test/fileTest2/gofile/pipe"
)

var f = flag.String("f", "test.txt", "input the file name")
var hostname = flag.String("s", "", "input server ip")
var port = flag.String("p", "", "input server port")

func init() {
	flag.Parse()
}

func main() {
	var ctrl pipe.Control
	if *hostname != "" {
		fmt.Println("send file")
		c := client.New(*hostname)
		ctrl = c
		bbytes := internal.Defaultbuffer.GetBytesbuffer(*f)
		ctrl.Write(bbytes)
	} else if *port != "" {
		fmt.Println("get file")
		s := &server.Server{
			Conn: make(chan net.Conn),
		}

		ctrl = s
		go s.Run(*port)
		for {
			fmt.Println("begin for")
			bb := ctrl.Read()
			internal.Defaultbuffer.PutBytesbufferToFile(bb.Bytes())
		}
	} else {
		fmt.Println("Please input gofile -h for help")
	}

}
