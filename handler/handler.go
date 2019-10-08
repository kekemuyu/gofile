package handler

import (
	"gofile/msg"
	"io"

	"github.com/donnie4w/go-logger/logger"
)

type Handler struct {
}

func HandleLoop(rwc io.ReadWriteCloser) {
	tmpb := make([]byte, 8)
	var message msg.Msg
	for {

		_, err := rwc.Read(tmpb)
		if err != nil {
			logger.Error("read head err:", err)
			continue
		}
		if message, err = message.Unpack(tmpb); err != nil {
			logger.Error("unpack msg err:", err)
			continue
		}

		message.Data = make([]byte, message.Datalen)
		_, err = rwc.Read(message.Data)
		if err != nil {
			logger.Error("read data err:", err)
			continue
		}

	}
}
