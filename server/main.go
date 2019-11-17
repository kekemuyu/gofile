package main

import (
	"gofile/config"
	"gofile/handler"
	_ "gofile/log"

	"gofile/com"
	"gofile/server/net"
	"runtime"

	"github.com/donnie4w/go-logger/logger"
)

func ConfigRuntime() {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	logger.Debug("Running with CPUs:", nuCPU)
}

func AppVersion() {
	logger.Debug("app-version:", config.Cfg.Section("").Key("app_ver").String())
}

func main() {
	portName := config.Cfg.Section("serial").Key("PortName").MustString("com1")
	baund := config.Cfg.Section("serial").Key("BaudRate").MustInt()
	defaultCom, err := com.New(portName, uint(baund))

	netport := config.Cfg.Section("server").Key("port").MustString(":9000")

	ConfigRuntime()
	AppVersion()
	if err != nil {
		logger.Error("打开串口错误：", err)
	} else {
		logger.Debug("打开串口：", portName)
		hlr := handler.Handler{
			Rwc:          defaultCom,
			Listch:       make(chan []string, 10),
			Uploadbodych: make(chan bool),
		}
		go hlr.HandleLoop() //start comm server
		logger.Debug(defaultCom)
	}
	logger.Debug("net serve:", netport)
	net.DefaultServer.Run(netport) //start net server

}
