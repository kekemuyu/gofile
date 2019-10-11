package filehandler

import (
	"os"
)

type Fileinfo struct {
	Name   string
	Parent string
	Isdrec byte
}

type Handler struct {
	Filehandler *os.File
	Name        string
	Size        int64
	Off         int64
}

var DefaultUpload, DefaultDownload *Handler

func (f *Fileinfo) pack() {
	// os.Create()
}
