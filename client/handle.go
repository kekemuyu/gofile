package main

import (
	"encoding/json"
	"fmt"
	"gofile/config"
	"gofile/handler"
	"gofile/msg"

	"os"

	log "github.com/donnie4w/go-logger/logger"
)

type FileHead struct {
	Name string
	Size int64
}

type Comtask struct {
	Name            string
	Size            int64
	Downloadfileoff int64
}

var DefaultComtask Comtask

func (c Comtask) ListServerpath() {
	message := msg.Msg{
		Id:      handler.Clist,
		Datalen: 0,
	}
	log.Debug(message)
	hlr.Sendmsg(message)
}

func (c Comtask) ListUppageServerpath() {
	message := msg.Msg{
		Id:      handler.Clistuppage,
		Datalen: 0,
	}
	log.Debug(message)
	hlr.Sendmsg(message)
}

//更新界面服务端目录
func (c Comtask) CListHandle(data []byte) {
	filemap := make(map[string][]string)
	err := json.Unmarshal(data, &filemap)
	if err != nil {
		log.Error(err)
		return
	}
	files := filemap["value"]

	jsStr := `$("#downfilesgroup").find("li").remove()`
	Defaultweb.UI.Eval(jsStr)
	for _, f := range files {
		log.Debug(f)

		jsStr := fmt.Sprintf(`$('#downfilesgroup').append("<li>%s</li>")`, f)
		Defaultweb.UI.Eval(jsStr)
	}
}

//更新上一界面服务端目录
func (c Comtask) CListUppageHandle(data []byte) {
	filemap := make(map[string][]string)
	err := json.Unmarshal(data, &filemap)
	if err != nil {
		log.Error(err)
		return
	}
	files := filemap["value"]

	jsStr := `$("#downfilesgroup").find("li").remove()`
	Defaultweb.UI.Eval(jsStr)
	for _, f := range files {
		log.Debug(f)

		jsStr := fmt.Sprintf(`$('#downfilesgroup').append("<li>%s</li>")`, f)
		Defaultweb.UI.Eval(jsStr)
	}
}

//向服务端发送head
func (c Comtask) DownloadHeadSend(name string) {
	hlr.Downname = name
	bs := []byte(name)
	message := msg.Msg{
		Id:      handler.Cdownloadhead,
		Datalen: uint32(len(bs)),
		Data:    bs,
	}

	hlr.Sendmsg(message)
}

func (c Comtask) CDownloadheadHandle(data []byte) {
	var fhead FileHead
	err := json.Unmarshal(data, &fhead)
	if err != nil {
		log.Error(err)
		return
	}
	c.Name = fhead.Name
	c.Size = fhead.Size
	c.Downloadfileoff = 0
	message := msg.Msg{
		Id:      handler.Cdownloadbody,
		Datalen: 0,
	}

	hlr.Sendmsg(message)
}

func (c Comtask) CDownloadbodyHandle(data []byte) {
	curpath := config.Cfg.Section("file").Key("clientpath").MustString(config.GetRootdir())
	path := curpath + `\` + c.Name
	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		log.Error(err)
		return
	}

	file.WriteAt(data, c.Downloadfileoff)
	c.Downloadfileoff += int64(len(data))
	if c.Downloadfileoff >= c.Size {
		return
	}
}

func (c Comtask) Upload(name string) {
	log.Debug(name)
	curpath := config.GetRootdir() + `\` + name

	file, err := os.Open(curpath)
	if err != nil {
		log.Error(err)
		return
	}
	fileinfo, err := os.Stat(curpath)
	if err != nil {
		log.Error(err)
		return
	}
	type filehead struct {
		Name string
		Size int64
	}

	fhead := filehead{
		Name: fileinfo.Name(),
		Size: fileinfo.Size(),
	}
	log.Debug(fhead)
	bs, _ := json.Marshal(fhead)
	message := msg.Msg{
		Id:      handler.Cuploadhead,
		Datalen: uint32(len(bs)),
		Data:    bs,
	}
	log.Debug(message)
	hlr.Sendmsg(message)

	blocksize := fhead.Size / 1024
	lastsize := fhead.Size % 1024
	log.Debug(blocksize, lastsize)
	if err != nil {
		log.Fatal(err)
	}
	outbytes := make([]byte, 1024)

	message = msg.Msg{
		Id:      handler.Cuploadbody,
		Datalen: 1024,
	}
	for i := int64(0); i < blocksize; {
		_, err = file.ReadAt(outbytes, i*1024)
		if err != nil {
			log.Error(err)
			return
		}
		message.Data = outbytes
		hlr.Sendmsg(message)
		<-hlr.Uploadbodych
	}
	n, _ := file.ReadAt(outbytes, blocksize*1024)
	if n > 0 {
		message.Datalen = uint32(lastsize)
		message.Data = outbytes[:lastsize]
		hlr.Sendmsg(message)
	}

}
