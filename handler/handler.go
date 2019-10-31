package handler

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"gofile/config"
	"gofile/filehandler"
	"gofile/msg"
	"io"
	"os"

	"github.com/donnie4w/go-logger/logger"
)

const (
	List uint32 = iota
	Uploadhead
	Uploadbody
	Downloadhead
	Downloadbody
)

type Handler struct {
	Rwc io.ReadWriteCloser
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
			logger.Error("read head err:", err)
			continue
		}

		if message, err = msg.Unpack(tmpb); err != nil {
			logger.Error("unpack msg err:", err)
			continue
		}

		message.Data = make([]byte, message.Datalen)
		n, err = h.Rwc.Read(message.Data)
		if n <= 0 {
			logger.Error("read data err:", err)
			continue
		}
		h.parseMsg(message)
	}
}

func (h *Handler) Send(data []byte) {
	h.Rwc.Write(data)
}

//procol magnage
func (h *Handler) parseMsg(msg msg.Msg) {
	switch msg.Id {
	case List:
		logger.Debug("df")
	case Uploadhead:
		logger.Debug("")
	case Uploadbody:
		logger.Debug("")
	case Downloadhead:
		logger.Debug("")
	case Downloadbody:
		logger.Debug("")
	}
}

func list(data []byte) {
	dir := string(data)
	if dir == "." { //list cur dir files,save to config json

		logger.Debug(".")

	}
	//list mou yi ge dir files
}

func uploadhead(data []byte) {
	type filehead struct {
		name string
		size int64
	}
	var fhead = filehead{}
	var err error
	if err = json.Unmarshal(data, &fhead); err != nil {
		logger.Error(err)
	}
	curpath := config.Cfg.Section("file").Key("path").MustString("/")
	if filehandler.DefaultUpload.Filehandler, err = os.Create(curpath + fhead.name); err != nil {
		logger.Error(err)
	}
	filehandler.DefaultUpload.Name = fhead.name
	filehandler.DefaultUpload.Size = fhead.size
	filehandler.DefaultUpload.Off = 0
}

func uploadbody(data []byte) {
	if filehandler.DefaultUpload.Filehandler == nil {
		return
	}
	filehandler.DefaultUpload.Filehandler.WriteAt(data, filehandler.DefaultUpload.Off)
	filehandler.DefaultUpload.Off += int64(len(data))
	if filehandler.DefaultUpload.Off >= filehandler.DefaultUpload.Size {
		filehandler.DefaultUpload.Filehandler.Close()
	}
}

func downloadhead(data []byte, rwc io.ReadWriteCloser) {
	name := string(data) //filename
	var err error
	if filehandler.DefaultDownload.Filehandler, err = os.Open(name); err != nil {
		logger.Error(err)
		return
	}

	var fileInfo os.FileInfo
	if fileInfo, err = filehandler.DefaultDownload.Filehandler.Stat(); err != nil {
		logger.Error(err)
		return
	}

	filehandler.DefaultDownload.Size = fileInfo.Size()
	var outmsg msg.Msg
	dataBuff := bytes.NewBuffer([]byte{})

	if err := binary.Write(dataBuff, binary.LittleEndian, filehandler.DefaultDownload.Size); err != nil {
		logger.Error(err)
		return
	}
	outmsg.Data = dataBuff.Bytes()
	outmsg.Id = Downloadhead
	outmsg.Datalen = uint32(len(outmsg.Data))

	outbytes, _ := msg.Pack(outmsg)
	rwc.Write(outbytes)
	filehandler.DefaultDownload.Off = 0
	filehandler.DefaultDownload.Blocksize = 1024
	filehandler.DefaultDownload.Blocknum = filehandler.DefaultDownload.Size / filehandler.DefaultDownload.Blocksize
	filehandler.DefaultDownload.Lastpacksize = filehandler.DefaultDownload.Size % filehandler.DefaultDownload.Blocksize
}

func downloadbody(rwc io.ReadWriteCloser) {
	var outmsg msg.Msg
	outmsg.Id = Downloadbody

	for i := 0; i < int(filehandler.DefaultDownload.Blocknum); i++ {
		_, err := filehandler.DefaultDownload.Filehandler.ReadAt(outmsg.Data, filehandler.DefaultDownload.Off)
		if err != nil {
			logger.Error(err)
			filehandler.DefaultDownload.Filehandler.Close()
			rwc.Close()
			return
		}
		filehandler.DefaultDownload.Off += filehandler.DefaultDownload.Blocksize
		outbytes, _ := msg.Pack(outmsg)
		_, err = rwc.Write(outbytes)
		if err != nil {
			rwc.Close()
			return
		}
	}
	if filehandler.DefaultDownload.Lastpacksize > 0 {
		filehandler.DefaultDownload.Filehandler.ReadAt(outmsg.Data, filehandler.DefaultDownload.Off)
		outbytes, _ := msg.Pack(outmsg)
		rwc.Write(outbytes)
		filehandler.DefaultDownload.Filehandler.Close()
		return //

	}
}
