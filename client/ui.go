package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"

	log "github.com/donnie4w/go-logger/logger"

	"github.com/zserge/lorca"
)

type Myweb struct {
	UI lorca.UI
}

var Defaultweb Myweb

func New(width, height int) Myweb {
	var myweb Myweb
	var err error
	myweb.UI, err = lorca.New("", "", width, height)
	if err != nil {
		log.Fatal(err)
	}

	return myweb
}

func (m *Myweb) Run() {

	ui := m.UI
	defer ui.Close()

	// A simple way to know when UI is ready (uses body.onload even in JS)
	ui.Bind("start", func() {
		log.Debug("UI is ready")
	})

	ui.Bind("opencom", Opencom)
	ui.Bind("browseclientpath", Browseclientpath)
	ui.Bind("browseclientuppage", Browseclientuppage)
	ui.Bind("upload", DefaultComtask.Upload)
	ui.Bind("download", DefaultComtask.DownloadHeadSend)
	ui.Bind("browseserverpath", DefaultComtask.ListServerpath)
	ui.Bind("browseserverpageup", DefaultComtask.ListUppageServerpath)

	//	ui.Bind("openSerial", New)
	//	ui.Bind("closeSerial", closeSerial)
	//	ui.Bind("send", send)
	//	ui.Bind("hexRecieveSet", hexRecieveSet)
	// Load HTML.
	// You may also use `data:text/html,<base64>` approach to load initial HTML,
	// e.g: ui.Load("data:text/html," + url.PathEscape(html))

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	go http.Serve(ln, http.FileServer(FS))
	ui.Load(fmt.Sprintf("http://%s", ln.Addr()))

	//	go recieve(ui)
	// You may use console.log to debug your JS code, it will be printed via
	// log.Println(). Also exceptions are printed in a similar manner.
	// ui.Eval(`
	// 	console.log("Hello, world!");
	// 	console.log('Multiple values:', [1, false, {"x":5}]);
	// `)

	// Wait until the interrupt signal arrives or browser window is closed
	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:
	case <-ui.Done():
	}

	log.Debug("exiting...")
}
