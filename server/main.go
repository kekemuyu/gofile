package main

import (
	"gofile/config"
	"gofile/handler"
	//_ "gofile/log"

	"gofile/com"
	"gofile/server/net"
	"runtime"
	"log"
	//"github.com/donnie4w/go-logger/logger"
)

func init(){
	log.SetFlags(log.Ldate|log.Lshortfile)
}

func ConfigRuntime() {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	log.Println("Running with CPUs:", nuCPU)
}

func AppVersion() {
	log.Println("app-version:", config.Cfg.Section("").Key("app_ver").String())
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
		log.Println("打开串口错误：", err)
	} else {
		log.Println("打开串口：", portName)

		hlr = handler.Handler{
			Rwc:    defaultCom,
			Listch: make(chan []string, 10),

			Shandler: &DefaultComtask,
		}

		go hlr.HandleLoop() //start comm server
		log.Println(defaultCom)
	}
	log.Println("net serve:", netport)
	net.DefaultServer.Run(netport) //start net server

}
