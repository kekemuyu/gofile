package handler

import (
	"gofile/msg"
	"gofile/protocol"
	"io"

	log "github.com/donnie4w/go-logger/logger"
)

const (
	Clist uint32 = iota
	Slist
	Clistuppage
	Slistuppage

	Cuploadhead
	Cuploadbody

	Cdownloadhead
	Cdownloadbody

	Sdownloadhead
	Sdownloadbody
)

type Handler struct {
	Rwc          io.ReadWriteCloser
	Listch       chan []string
	Downname     string
	Downoff      int64
	Downsize     int64
	Uploadbodych chan bool
	Chandler     protocol.CHandler
	Shandler     protocol.SHandler
}

func (h *Handler) HandleLoop() {
	tmpb := make([]byte, 8)
	var message msg.Msg
	for {

		n, err := h.Rwc.Read(tmpb)
		if n <= 0 {
			continue
		}
		if err != nil {
			log.Error("read head err:", err)
			continue
		}

		if message, err = msg.Unpack(tmpb); err != nil {
			log.Error("unpack msg err:", err)
			continue
		}

		log.Debug(message)
		if message.Datalen > 0 {

			message.Data = make([]byte, message.Datalen)
			n, err = h.Rwc.Read(message.Data)
			if n <= 0 {
				log.Error("read data err:", err)
				continue
			}
			h.parseMsg(message)
		} else {
			message.Data = make([]byte, 1)
			h.parseMsg(message)
		}
	}
}

func (h *Handler) Send(data []byte) {
	h.Rwc.Write(data)
}

func (h *Handler) Sendmsg(message msg.Msg) {
	bs, err := msg.Pack(message)
	if err != nil {
		log.Error("sendmsg err:", err)
		return
	}
	h.Send(bs)
}

//procol magnage
func (h *Handler) parseMsg(msg msg.Msg) {
	switch msg.Id {

	case Clist: //客户端发送浏览,服务端处理
		log.Debug("Clist")

		h.Shandler.SListHandle(msg.Data)

	case Slist: //服务端发送浏览结果
		log.Debug("Slist")
		h.Chandler.CListHandle(msg.Data)
	case Clistuppage:
		log.Debug("Clistuppage")
		h.Shandler.SListUppageHandle(msg.Data)
	case Slistuppage:
		log.Debug("Slistuppage")
		h.Chandler.CListUppageHandle(msg.Data)

	case Cuploadhead:
		log.Debug("Cuploadhead")
		h.Shandler.SUploadheadHandle(msg.Data)
	case Cuploadbody:
		log.Debug("Cuploadbody")
		h.Shandler.SUploadbodyHandle(msg.Data)

	case Cdownloadhead:
		log.Debug("Cdownloadhead")
		h.Shandler.SDownloadheadHandle(msg.Data)
	case Cdownloadbody:
		log.Debug("Cdownloadbody")
		h.Shandler.SDownloadbodyHandle(msg.Data)
	case Sdownloadhead:

		log.Debug("Sdownloadhead")
		h.Chandler.CDownloadheadHandle(msg.Data)
	case Sdownloadbody:
		log.Debug("Sdownloadbody")
		h.Chandler.CDownloadbodyHandle(msg.Data)
	}
}
