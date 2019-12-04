package handler

import (
	"gofile/msg"
	"gofile/protocol"
	"io"
	"log"
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
		if n <= 0 {
			continue
		}
		if err != nil {
			log.Println("read head err:", err)
			continue
		}

		log.Println(tmpb)
		if message, err = msg.Unpack(tmpb); err != nil {
			log.Println("unpack msg err:", err)
			continue
		}

		log.Println(message)
		if message.Datalen > 0 {

			message.Data = make([]byte, message.Datalen)

			dataP := 0
			for dataP < int(message.Datalen) {
				n, err = h.Rwc.Read(message.Data[dataP:])
				if n > 0 {
					dataP += n
				}
				if err != nil {
					log.Println("read data err:", err)
					break
				}
			}

			log.Println(message)
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
		log.Println("sendmsg err:", err)
		return
	}
	h.Send(bs)
}

//procol magnage
func (h *Handler) parseMsg(msg msg.Msg) {
	switch msg.Id {

	case Clist: //客户端发送浏览,服务端处理
		log.Println("Clist")

		h.Shandler.SListHandle(msg.Data)

	case Slist: //服务端发送浏览结果
		log.Println("Slist")
		h.Chandler.CListHandle(msg.Data)
	case Clistuppage:
		log.Println("Clistuppage")
		h.Shandler.SListUppageHandle(msg.Data)
	case Slistuppage:
		log.Println("Slistuppage")
		h.Chandler.CListUppageHandle(msg.Data)

	case Cuploadhead:
		log.Println("Cuploadhead")
		h.Shandler.SUploadheadHandle(msg.Data)
	case Cuploadbody:
		log.Println("Cuploadbody")
		h.Shandler.SUploadbodyHandle(msg.Data)
	case Suploadbody_nextpack: //服务端收到数据后，发送响应到客户端
		log.Println("Suploadbody")
		h.Chandler.CUploadbodyNextpackHandle(msg.Data)

	case Cdownloadhead:
		log.Println("Cdownloadhead")
		h.Shandler.SDownloadheadHandle(msg.Data)
	case Cdownloadbody:
		log.Println("Cdownloadbody")
		go h.Shandler.SDownloadbodyHandle(msg.Data)
	case Sdownloadhead:

		log.Println("Sdownloadhead")
		h.Chandler.CDownloadheadHandle(msg.Data)
	case Sdownloadbody:
		log.Println("Sdownloadbody")
		h.Chandler.CDownloadbodyHandle(msg.Data)
	case Cdownloadbody_nextpack:
		log.Println("Cdownloadbody_nextpack")
		h.Shandler.SDownloadbodyNextpackHandle(msg.Data)
	case Clistdisk:
		h.Shandler.SListdisk(msg.Data)
	}

}
