package filehandler

import (
	"os"
)

type FileHead struct {
	Name string
	Size int64
}

type Handler struct {
	Filehandler  *os.File
	Name         string
	Size         int64
	Off          int64
	Blocksize    int64
	Blocknum     int64
	Lastpacksize int64
}

var DefaultUpload, DefaultDownload Handler
