package internal

import (
	"gofile/msg"

	"github.com/donnie4w/go-logger/logger"
)

type IO interface {
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
	ReadHandle(msg msg.Msg)
}

func HandleLoop(io IO) {
	tmpb := make([]byte, 8)
	var message msg.Msg
	for {

		_, err := io.Read(tmpb)
		if err != nil {
			logger.Error("read head err:", err)
			continue
		}
		if message, err = message.Unpack(tmpb); err != nil {
			logger.Error("unpack msg err:", err)
			continue
		}

		message.Data = make([]byte, message.Datalen)
		_, err = io.Read(message.Data)
		if err != nil {
			logger.Error("read data err:", err)
			continue
		}
		go io.ReadHandle(message)

	}
}
