package main

import (
	"gofile/config"
	"gofile/internal"
	_ "gofile/log"

	"gofile/server/com"
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
		logger.Debug(defaultCom)
	}

	net.DefaultServer.Run(netport)
	internal.HandleLoop(net.DefaultServer)

}
