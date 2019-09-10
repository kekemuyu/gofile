package com

import (
	"gofile/internal"
	"io"

	"github.com/donnie4w/go-logger/logger"
	"github.com/jacobsa/go-serial/serial"
)

type Com struct {
	IOcom   io.ReadWriteCloser
	ReadCh  chan []byte
	WriteCh chan []byte
}

func New(portnum string, baudrate uint) (*Com, error) {
	opt := serial.OpenOptions{
		PortName:        portnum,
		BaudRate:        baudrate,
		DataBits:        8,
		StopBits:        1,
		ParityMode:      serial.PARITY_NONE,
		MinimumReadSize: 50,
	}

	var err error
	var tempIO io.ReadWriteCloser

	tempIO, err = serial.Open(opt)

	return &Com{
		IOcom:   tempIO,
		ReadCh:  make(chan []byte, 10),
		WriteCh: make(chan []byte, 10),
	}, err
}

func (c *Com) Run(io internal.IO) {
	buf := make([]byte, 1024)
	go c.Comwrite() //串口发送
	for {
		cnt, err := c.IOcom.Read(buf)
		if err != nil {
			logger.Error("com read err:", err)
			continue
		}
		if cnt > 0 {
			io.Read(buf[:cnt])
		}
	}
}

func (c *Com) Comwrite() {
	for {
		select {
		case writech := <-c.WriteCh:
			_, err := c.IOcom.Write(writech)
			if err != nil {
				logger.Error("Chansend err:", err)
			}
		default:
		}
	}
}
