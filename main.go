package main

import (
	"flag"
	"fmt"
	"net"

	"test/fileTest2/gofile/internal"
	"test/fileTest2/gofile/networker/client"
	"test/fileTest2/gofile/networker/server"
	"test/fileTest2/gofile/pipe"
	"test/fileTest2/gofile/serialworker"
)

var file = flag.String("f", "test.txt", "input the file name")
var hostname = flag.String("s", "", "input server ip")
var port = flag.String("p", "", "input server port")
var com = flag.String("c", "", "input com port")
var mode = flag.String("m", "", "input r or s for recieve or send ")

func init() {
	flag.Parse()
}

func main() {
	var ctrl pipe.Control
	if *hostname != "" {
		fmt.Println("send file")
		c := client.New(*hostname)
		ctrl = c
		if *file == "" {
			fmt.Println("input the file name with -f ")
			return
		}
		bbytes := internal.Defaultbuffer.GetBytesbuffer(*file)
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
	} else if (*com != "") && (*mode != "") {
		fmt.Println("opened com port is:", *com)
		w := serialworker.New(*com)
		ctrl = w
		if *mode == "s" {
			if *file == "" {
				fmt.Println("input the file name with -f ")
				return
			}
			bbytes := internal.Defaultbuffer.GetBytesbuffer(*file)
			ctrl.Write(bbytes)
		} else if *mode == "r" {
			for {
				bb := ctrl.Read()
				internal.Defaultbuffer.PutBytesbufferToFile(bb.Bytes())

			}
		} else {
			fmt.Println("input the serial mode  with -r or -s ")
		}

	} else {
		fmt.Println("Please input gofile -h for help")
	}

}
