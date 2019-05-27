package main

import (
	"flag"
	"test/fileTest2/gofile/internal"
)

var f = flag.String("f", "1.txt", "input the file name")
var hostname = flag.String("s", "127.0.0.1:5000", "input server ip")

func init() {
	flag.Parse()
}

func main() {
	//	internal.Defaultbuffer.Run(*f)   //
}
