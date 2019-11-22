package main

import (
	"encoding/json"
	"gofile/config"
	"gofile/handler"
	"gofile/msg"
	"gofile/util"

	"io/ioutil"

	"os"

	log "github.com/donnie4w/go-logger/logger"
)

type FileHead struct {
	Name string
	Size int64
}

type Comtask struct {
	Serverpath    string
	Name          string
	Size          int64
	Handler       *os.File
	Uploadfileoff int64
}

var DefaultComtask Comtask

//将当前目录文件发送给客户端
func (c Comtask) SListHandle(data []byte) {
	log.Debug(data)
	curpath := config.Cfg.Section("file").Key("serverpath").MustString(config.GetRootdir())
	path := curpath
	if len(data) != 1 {
		path = curpath + `\` + string(data)
	}

	log.Debug(path)

	files, _ := ioutil.ReadDir(path)

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
			Id:      handler.Slist,
			Datalen: uint32(len(bs)),
			Data:    bs,
		}
		hlr.Sendmsg(message)
		config.Cfg.Section("file").Key("serverpath").SetValue(path)
		config.Save()
	}
}

//将上一页目录文件发送给客户端
func (c Comtask) SListUppageHandle(data []byte) {

	curpath := config.Cfg.Section("file").Key("serverpath").MustString(config.GetRootdir())
	curpath = util.GetParentDirectory(curpath)
	path := curpath + `\` + string(data)

	log.Debug(path)

	files, _ := ioutil.ReadDir(path)

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
			Id:      handler.Slistuppage,
			Datalen: uint32(len(bs)),
			Data:    bs,
		}
		hlr.Sendmsg(message)
		config.Cfg.Section("file").Key("serverpath").SetValue(path)
		config.Save()
	}
}

func (c Comtask) SUploadheadHandle(data []byte) {

	var fhead FileHead
	var err error
	if err = json.Unmarshal(data, &fhead); err != nil {
		log.Error(err)
	}
	log.Debug(fhead)
	curpath := config.Cfg.Section("file").Key("serverpath").MustString(config.GetRootdir())
	path := curpath + `\` + fhead.Name
	if c.Handler, err = os.Create(path); err != nil {
		log.Error(err)
	}
	c.Name = fhead.Name
	c.Size = fhead.Size
	c.Uploadfileoff = 0
}

func (c Comtask) SUploadbodyHandle(data []byte) {
	c.Handler.WriteAt(data, c.Uploadfileoff)
	c.Uploadfileoff += int64(len(data))
	if c.Uploadfileoff >= c.Size {
		c.Handler.Close()
		return
	}
}

func (c Comtask) SDownloadheadHandle(data []byte) {
	curpath := config.Cfg.Section("file").Key("serverpath").MustString(config.GetRootdir())
	path := curpath + `\` + string(data)
	var err error
	c.Handler, err = os.Open(path)
	if err != nil {
		log.Error(err)
		return
	}
	fileinfo, err := os.Stat(path)
	if err != nil {
		log.Error(err)
		return
	}

	fhead := FileHead{
		Name: fileinfo.Name(),
		Size: fileinfo.Size(),
	}
	c.Name = fhead.Name
	c.Size = fhead.Size
	log.Debug(fhead)
	bs, _ := json.Marshal(fhead)
	message := msg.Msg{
		Id:      handler.Sdownloadhead,
		Datalen: uint32(len(bs)),
		Data:    bs,
	}
	log.Debug(message)
	hlr.Sendmsg(message)
}

func (c Comtask) SDownloadbodyHandle(data []byte) {
	blocksize := c.Size / 1024
	lastsize := c.Size % 1024
	log.Debug(blocksize, lastsize)

	outbytes := make([]byte, 1024)

	message := msg.Msg{
		Id:      handler.Sdownloadbody,
		Datalen: 1024,
	}
	for i := int64(0); i < blocksize; {
		_, err := c.Handler.ReadAt(outbytes, i*1024)
		if err != nil {
			log.Error(err)
			return
		}
		message.Data = outbytes
		hlr.Sendmsg(message)
	}
	if lastsize > 0 {
		n, _ := c.Handler.ReadAt(outbytes, blocksize*1024)
		if n > 0 {
			message.Datalen = uint32(lastsize)
			message.Data = outbytes[:lastsize]
			hlr.Sendmsg(message)

		}
	}
	c.Handler.Close()
}
