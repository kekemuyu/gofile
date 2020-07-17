package main

import (
	"gofile/config"
	"gofile/handler"

	_ "gofile/log"
	"runtime"

	"gofile/com"
	"gofile/server/net"
	// log "github.com/donnie4w/go-logger/logger"
)

func ConfigRuntime() {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	// log.Debug("Running with CPUs:", nuCPU)
}

func AppVersion() {
	// log.Debug("app-version:", config.Cfg.Section("").Key("app_ver").String())
}

var hlr handler.Handler

func main() {
	portName := config.Cfg.Section("serial").Key("PortName").MustString("com1")
	baund := config.Cfg.Section("serial").Key("BaudRate").MustInt()
	defaultCom, err := com.New(portName, uint(baund))

	netport := config.Cfg.Section("server").Key("port").MustString(":9000")

	ConfigRuntime()
	AppVersion()
	if err != nil {
		// log.Error("打开串口错误：", err)
	} else {
		// log.Debug("打开串口：", portName)

		hlr = handler.Handler{
			Rwc:    defaultCom,
			Listch: make(chan []string, 10),

			Shandler: &DefaultComtask,
		}

		go hlr.HandleLoop() //start comm server
		// log.Debug(defaultCom)
	}
	// log.Debug("net serve:", netport)
	net.DefaultServer.Run(netport) //start net server

}
