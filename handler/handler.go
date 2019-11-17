package handler

import (
	"encoding/json"
	"gofile/config"
	"gofile/filehandler"
	"gofile/msg"
	"io"
	"io/ioutil"
	"os"

	log "github.com/donnie4w/go-logger/logger"
)

const (
	List uint32 = iota
	Relist
	Uploadhead
	Uploadbody
	Reuploadbody
	Download
	Redownload
)

type Handler struct {
	Rwc          io.ReadWriteCloser
	Listch       chan []string
	Downname     string
	Downoff      int64
	Downsize     int64
	Uploadbodych chan bool
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

		if message.Datalen <= 0 {
			h.parseMsg(message)
			continue
		}
		message.Data = make([]byte, message.Datalen)
		n, err = h.Rwc.Read(message.Data)
		if n <= 0 {
			log.Error("read data err:", err)
			continue
		}
		h.parseMsg(message)
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
	case List:
		log.Debug("list")
		h.list("1")
	case Relist:
		log.Debug("Relist")
		h.Relist(msg.Data)
	case Uploadhead:

		log.Debug("Uploadhead")
		uploadhead(msg.Data)
	case Uploadbody:
		log.Debug("Uploadbody")
		h.uploadbody(msg.Data)
	case Reuploadbody:
		log.Debug("Reuploadbody")
		h.Reuploadbody()
	case Download:
		log.Debug("Download")
		h.download(msg.Data)
	case Redownload:
		log.Debug("Redownload")
		h.Redownload(msg.Data)
	}
}

func (h *Handler) list(dirname string) {
	curpath := config.GetRootdir()

	files, _ := ioutil.ReadDir(curpath)

	var filenames []string
	for _, f := range files {
		filenames = append(filenames, f.Name())
		log.Debug(f.Name())

	}

	filemap := make(map[string][]string)
	filemap["value"] = filenames

	bs, err := json.Marshal(filemap)
	if err != nil {
		log.Error(err)
		return
	}

	if len(filenames) > 0 {
		message := msg.Msg{
			Id:      Relist,
			Datalen: uint32(len(bs)),
			Data:    bs,
		}
		h.Sendmsg(message)
	}

}

func (h *Handler) Relist(data []byte) {
	filemap := make(map[string][]string)
	err := json.Unmarshal(data, &filemap)
	if err != nil {
		log.Error(err)
		return
	}
	log.Debug(filemap["value"])
	h.Listch <- filemap["value"] //send server local dir files to client
}
func uploadhead(data []byte) {
	type filehead struct {
		Name string
		Size int64
	}
	var fhead = filehead{}
	var err error
	if err = json.Unmarshal(data, &fhead); err != nil {
		log.Error(err)
	}
	log.Debug(fhead)
	// curpath := config.GetRootdir() + `\`
	log.Debug(fhead.Name)
	if filehandler.DefaultUpload.Filehandler, err = os.Create(fhead.Name); err != nil {
		log.Error(err)
	}
	filehandler.DefaultUpload.Name = fhead.Name
	filehandler.DefaultUpload.Size = fhead.Size
	filehandler.DefaultUpload.Off = 0
}

func (h *Handler) uploadbody(data []byte) {
	if filehandler.DefaultUpload.Filehandler == nil {
		return
	}
	filehandler.DefaultUpload.Filehandler.WriteAt(data, filehandler.DefaultUpload.Off)
	filehandler.DefaultUpload.Off += int64(len(data))
	if filehandler.DefaultUpload.Off >= filehandler.DefaultUpload.Size {
		log.Debug("uploadbody comlete")
		filehandler.DefaultUpload.Filehandler.Close()
		return
	}
	message := msg.Msg{
		Id:      Reuploadbody,
		Datalen: 0,
	}
	h.Sendmsg(message)
}

func (h *Handler) Reuploadbody() {
	h.Uploadbodych <- true
}

func (h *Handler) download(data []byte) {
	name := string(data) //filename
	var err error
	var file *os.File
	if file, err = os.Open(name); err != nil {
		log.Error(err)
		return
	}

	var fileInfo os.FileInfo
	if fileInfo, err = os.Stat(name); err != nil {
		log.Error(err)
		return
	}

	filesize := fileInfo.Size()
	blocksize := filesize / 1024
	lastsize := filesize % 1024
	log.Debug(blocksize, lastsize)

	outbytes := make([]byte, 1024)
	h.Downoff = 0
	h.Downsize = filesize
	message := msg.Msg{
		Id:      Redownload,
		Datalen: 1024,
	}
	for i := int64(0); i < blocksize; {
		_, err = file.ReadAt(outbytes, i*1024)
		if err != nil {
			log.Error(err)
			return
		}
		message.Data = outbytes
		h.Sendmsg(message)

	}
	n, _ := file.ReadAt(outbytes, blocksize*1024)
	log.Debug(n)
	if n > 0 {

		message.Datalen = uint32(lastsize)
		message.Data = outbytes[:lastsize]
		log.Debug(message)
		h.Sendmsg(message)
	}
}

func (h *Handler) Redownload(data []byte) {
	log.Debug(h.Downname)
	file, err := os.Create(h.Downname)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
		return
	}

	file.WriteAt(data, h.Downoff)
	h.Downoff += int64(len(data))
	if h.Downoff >= h.Downsize {
		return
	}
}
