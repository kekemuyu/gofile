package handler

import (
	"io"

	"github.com/kekemuyu/gofile/msg"
	"github.com/kekemuyu/gofile/protocol"
	//log "github.com/donnie4w/go-logger/logger"
)

const (
	Clist uint32 = iota
	Slist
	Clistuppage
	Slistuppage

	Cuploadhead
	Cuploadbody
	Suploadbody_nextpack

	Cdownloadhead
	Cdownloadbody
	Cdownloadbody_nextpack

	Sdownloadhead
	Sdownloadbody

	//关于驱动器的操作
	Clistdisk
	Slistdisk
)

type Handler struct {
	Rwc      io.ReadWriteCloser
	Listch   chan []string
	Downname string
	Downoff  int64
	Downsize int64

	Chandler protocol.CHandler
	Shandler protocol.SHandler
}

func (h *Handler) HandleLoop() {
	tmpb := make([]byte, 8)
	var message msg.Msg
	for {
		n, err := h.Rwc.Read(tmpb)
		if err != nil {
			// log.Error("read head err:", err)
			return
		}
		if n <= 0 {
			continue
		}

		// log.Debug(tmpb)
		if message, err = msg.Unpack(tmpb); err != nil {
			// log.Error("unpack msg err:", err)
			continue
		}

		// log.Debug(message)
		if message.Datalen > 0 {

			message.Data = make([]byte, message.Datalen)

			dataP := 0
			for dataP < int(message.Datalen) {
				n, err = h.Rwc.Read(message.Data[dataP:])
				if err != nil {
					// log.Error("read data err:", err)
					return
				}
				if n > 0 {
					dataP += n
				}

			}

			// log.Debug(message)
			go h.parseMsg(message)
		} else {
			message.Data = make([]byte, 1)
			go h.parseMsg(message)
		}
	}
}

func (h *Handler) Send(data []byte) {
	h.Rwc.Write(data)
}

func (h *Handler) Sendmsg(message msg.Msg) {
	bs, err := msg.Pack(message)
	if err != nil {
		// log.Error("sendmsg err:", err)
		return
	}
	h.Send(bs)
}

//procol magnage
func (h *Handler) parseMsg(msg msg.Msg) {
	switch msg.Id {

	case Clist: //客户端发送浏览,服务端处理
		// log.Debug("Clist")

		h.Shandler.SListHandle(msg.Data)

	case Slist: //服务端发送浏览结果
		// log.Debug("Slist")
		h.Chandler.CListHandle(msg.Data)
	case Clistuppage:
		// log.Debug("Clistuppage")
		h.Shandler.SListUppageHandle(msg.Data)
	case Slistuppage:
		// log.Debug("Slistuppage")
		h.Chandler.CListUppageHandle(msg.Data)

	case Cuploadhead:
		// log.Debug("Cuploadhead")
		h.Shandler.SUploadheadHandle(msg.Data)
	case Cuploadbody:
		// log.Debug("Cuploadbody")
		h.Shandler.SUploadbodyHandle(msg.Data)
	case Suploadbody_nextpack: //服务端收到数据后，发送响应到客户端
		// log.Debug("Suploadbody")
		h.Chandler.CUploadbodyNextpackHandle(msg.Data)

	case Cdownloadhead:
		// log.Debug("Cdownloadhead")
		h.Shandler.SDownloadheadHandle(msg.Data)
	case Cdownloadbody:
		// log.Debug("Cdownloadbody")
		go h.Shandler.SDownloadbodyHandle(msg.Data)
	case Sdownloadhead:

		// log.Debug("Sdownloadhead")
		h.Chandler.CDownloadheadHandle(msg.Data)
	case Sdownloadbody:
		// log.Debug("Sdownloadbody")
		h.Chandler.CDownloadbodyHandle(msg.Data)
	case Cdownloadbody_nextpack:
		// log.Debug("Cdownloadbody_nextpack")
		h.Shandler.SDownloadbodyNextpackHandle(msg.Data)
	case Clistdisk:
		h.Shandler.SListdisk(msg.Data)
	}

}
