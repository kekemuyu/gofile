package handler

import (
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

//procol magnage
func parseMsg(msg msg.Msg) {
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

func downloadhead(data []byte) {

}
func downloadbody(data []byte) {

}
