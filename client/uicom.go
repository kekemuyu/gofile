package main

import (
	"gofile/com"
	"gofile/handler"
	"gofile/msg"

	log "github.com/donnie4w/go-logger/logger"
)

var hlr handler.Handler

func Opencom(comnum string, baudrate int) bool {
	irw, err := com.New(comnum, uint(baudrate))
	if err != nil {
		return false
	}
	hlr = handler.Handler{Rwc: irw}
	log.Debug("Opencom:", hlr)
	go hlr.HandleLoop()
	return true
}

func Sendmsg(message msg.Msg) {
	bs, err := msg.Pack(message)
	if err != nil {
		log.Error("sendmsg err:", err)
		return
	}
	hlr.Send(bs)
}
