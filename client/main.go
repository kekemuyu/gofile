//go:generate go run -tags generate gen.go
package main

import (
	_ "gofile/log"
)

// "github.com/donnie4w/go-logger/logger"

func main() {
	// conn, err := net.Dial("tcp", "127.0.0.1:9000")
	// if err != nil {
	// 	panic(conn)
	// }
	// logger.Debug(conn, err)
	// var msg = msg.Msg{
	// 	Id: 0x00000001,
	// }
	// msg.Data = []byte("test.txt")
	// msg.Datalen = uint32(len(msg.Data))
	// logger.Debug(msg)
	// b, err := msg.Pack(msg)
	// logger.Debug(b)
	// conn.Write(b)

	Defaultweb = New(800, 600)
	Defaultweb.Run()
}
