package main

import (
	"encoding/json"
	"fmt"
	"strconv"

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
	Name                  string
	Size                  int64
	Downloadfileoff       int64
	Uploadbody_nextpackch chan bool
	Filehandler           *os.File
	Progress              Tprocess
}

var DefaultComtask = Comtask{
	Uploadbody_nextpackch: make(chan bool),
}

func (c *Comtask) ListServerpath(serverpathname string) {
	bs := []byte(serverpathname)
	message := msg.Msg{
		Id:      handler.Clist,
		Datalen: uint32(len(bs)),
		Data:    bs,
	}
	log.Debug(message)
	hlr.Sendmsg(message)
}

func (c *Comtask) ListUppageServerpath() {
	message := msg.Msg{
		Id:      handler.Clistuppage,
		Datalen: 2,
		Data:    []byte("up"),
	}
	log.Debug(message)
	hlr.Sendmsg(message)
}

//更新界面服务端目录
func (c *Comtask) CListHandle(data []byte) {
	filemap := make(map[string][]string)
	err := json.Unmarshal(data, &filemap)
	if err != nil {
		log.Error(err)
		return
	}

	log.Debug(filemap)
	files := filemap["value"]
	jsStr := fmt.Sprintf(`$('#serverpath').val("%s");$("#serverfiles").find("li").remove()`, files[0])
	Defaultweb.UI.Eval(jsStr)
	for _, f := range files[1:] {
		log.Debug(f)

		jsStr := fmt.Sprintf(`$('#serverfiles').append("<li>%s</li>")`, f)
		Defaultweb.UI.Eval(jsStr)
	}
}

//更新上一界面服务端目录
func (c *Comtask) CListUppageHandle(data []byte) {
	filemap := make(map[string][]string)
	err := json.Unmarshal(data, &filemap)
	if err != nil {
		log.Error(err)
		return
	}
	log.Debug(filemap)
	files := filemap["value"]

	jsStr := fmt.Sprintf(`$('#serverpath').val("%s");$("#serverfiles").find("li").remove()`, files[0])

	Defaultweb.UI.Eval(jsStr)
	for _, f := range files[1:] {
		log.Debug(f)

		jsStr := fmt.Sprintf(`$('#serverfiles').append("<li>%s</li>")`, f)
		Defaultweb.UI.Eval(jsStr)
	}
}

//向服务端发送head
func (c *Comtask) DownloadHeadSend(name string) {
	hlr.Downname = name
	bs := []byte(name)
	message := msg.Msg{
		Id:      handler.Cdownloadhead,
		Datalen: uint32(len(bs)),
		Data:    bs,
	}
	log.Debug(message)
	hlr.Sendmsg(message)
}

func (c *Comtask) CDownloadheadHandle(data []byte) {
	var fhead FileHead
	err := json.Unmarshal(data, &fhead)
	if err != nil {
		log.Error(err)
		return
	}
	c.Name = fhead.Name
	c.Size = fhead.Size
	c.Downloadfileoff = 0
	log.Debug(c.Name, c.Size)

	bs := []byte("downloadbody")
	message := msg.Msg{
		Id:      handler.Cdownloadbody,
		Datalen: uint32(len(bs)),
		Data:    bs,
	}

	hlr.Sendmsg(message)

	curpath := config.Cfg.Section("file1").Key("clientpath").MustString(config.GetRootdir())
	path := curpath + `/` + c.Name
	c.Filehandler, err = os.Create(path)
	if err != nil {
		log.Error(err)
		return
	}
}

func (c *Comtask) CDownloadbodyHandle(data []byte) {
	if c.Filehandler == nil {
		return
	}

	log.Debug(c.Downloadfileoff, c.Size)
	c.Filehandler.WriteAt(data, c.Downloadfileoff)
	c.Downloadfileoff += int64(len(data))
	c.Progress.Filename = c.Name
	c.Progress.Filesize = c.Size
	c.Progress.Runsize = c.Downloadfileoff

	go Disprocessbar(c.Progress) //进度条更新
	if c.Downloadfileoff >= c.Size {
		c.Filehandler.Close()
		log.Debug("download complete")
		return
	}

	bs := []byte("ok")
	message := msg.Msg{
		Id:      handler.Cdownloadbody_nextpack,
		Datalen: uint32(len(bs)),
		Data:    bs,
	}
	log.Debug(message)
	hlr.Sendmsg(message)
}

func (c *Comtask) Upload(name string) {
	log.Debug(name)
	curpath := config.Cfg.Section("file1").Key("clientpath").MustString(config.GetRootdir()) + `/` + name

	file, err := os.Open(curpath)
	defer file.Close()
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
	c.Progress.Filename = fhead.Name
	c.Progress.Filesize = fhead.Size
	for i := int64(0); i < blocksize; i++ {
		_, err = file.ReadAt(outbytes, i*1024)
		if err != nil {
			log.Error(err)
			return
		}
		message.Data = outbytes
		hlr.Sendmsg(message)
		<-c.Uploadbody_nextpackch
		// log.Debug(i)

		c.Progress.Runsize = i * 1024

		go Disprocessbar(c.Progress) //进度条更新
	}

	n, _ := file.ReadAt(outbytes, blocksize*1024)
	if n > 0 {
		message.Datalen = uint32(lastsize)
		message.Data = outbytes[:lastsize]
		hlr.Sendmsg(message)
	}
	c.Progress.Runsize = fhead.Size

	go Disprocessbar(c.Progress) //进度条更新

}

//收到服务器响应
func (c *Comtask) CUploadbodyNextpackHandle(data []byte) {
	c.Uploadbody_nextpackch <- true
}

func (c *Comtask) GetServerDisk() {
	bs := []byte("clistdisk")
	message := msg.Msg{
		Id:      handler.Clistdisk,
		Datalen: uint32(len(bs)),
		Data:    bs,
	}
	log.Debug(message)
	hlr.Sendmsg(message)
}

func (c *Comtask) CListdisk(data []byte) {
	var dinfo []string
	err := json.Unmarshal(data, &dinfo)
	if err != nil {
		return
	}
	jsStr := `$('#serverDisk').find("option").remove()`
	Defaultweb.UI.Eval(jsStr)
	if len(dinfo) > 0 {
		for k, v := range dinfo {
			jsStr := fmt.Sprintf(`$('#serverDisk').append("<option value=%s>%s</option>")`, strconv.Itoa(k), v)
			Defaultweb.UI.Eval(jsStr)
		}
	}
}
