package main

import (
	"encoding/json"
	"fmt"
	"gofile/com"
	"gofile/config"
	"gofile/handler"
	"gofile/msg"
	"io/ioutil"
	"os"

	// "path/filepath"

	log "github.com/donnie4w/go-logger/logger"
)

var hlr handler.Handler

func Opencom(comnum string, baudrate int) bool {
	irw, err := com.New(comnum, uint(baudrate))
	if err != nil {
		return false
	}
	hlr = handler.Handler{
		Rwc:          irw,
		Listch:       make(chan []string, 10),
		Uploadbodych: make(chan bool),
	}
	log.Debug("Opencom:", hlr)
	go hlr.HandleLoop()
	go Run()
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

func Browsecurpath() {
	curpath := config.GetRootdir()
	// files, _ := filepath.Glob(curpath)

	files, _ := ioutil.ReadDir(curpath)
	jsStr1 := `$("#filesgroup").find("li").remove()`
	Defaultweb.UI.Eval(jsStr1)
	for _, f := range files {
		log.Debug(f.Name())

		jsStr := fmt.Sprintf(`$('#filesgroup').append("<li>%s</li>")`, f.Name())
		Defaultweb.UI.Eval(jsStr)
	}
}

func Upload(name string) {
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
		Id:      handler.Uploadhead,
		Datalen: uint32(len(bs)),
		Data:    bs,
	}
	log.Debug(message)
	Sendmsg(message)

	blocksize := fhead.Size / 1024
	lastsize := fhead.Size % 1024
	log.Debug(blocksize, lastsize)
	if err != nil {
		log.Fatal(err)
	}
	outbytes := make([]byte, 1024)

	message = msg.Msg{
		Id:      handler.Uploadbody,
		Datalen: 1024,
	}
	for i := int64(0); i < blocksize; {
		_, err = file.ReadAt(outbytes, i*1024)
		if err != nil {
			log.Error(err)
			return
		}
		message.Data = outbytes
		Sendmsg(message)
		<-hlr.Uploadbodych
	}
	n, _ := file.ReadAt(outbytes, blocksize*1024)
	if n > 0 {
		message.Datalen = uint32(lastsize)
		message.Data = outbytes[:lastsize]
		Sendmsg(message)
	}

}

//send list command to server
func Browsedowncurpath() {
	message := msg.Msg{
		Id:      handler.List,
		Datalen: 0,
	}

	Sendmsg(message)
}

func Download(name string) {
	hlr.Downname = name
	bs := []byte(name)
	message := msg.Msg{
		Id:      handler.Download,
		Datalen: uint32(len(bs)),
		Data:    bs,
	}

	Sendmsg(message)
}

func Run() {
	select {
	case listdir := <-hlr.Listch:
		log.Debug(listdir)
		jsStr1 := `$("#downfilesgroup").find("li").remove()`
		Defaultweb.UI.Eval(jsStr1)
		for _, f := range listdir {
			log.Debug(f)

			jsStr := fmt.Sprintf(`$('#downfilesgroup').append("<li>%s</li>")`, f)
			Defaultweb.UI.Eval(jsStr)
		}

	}

}
