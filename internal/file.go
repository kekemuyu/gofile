package internal

import (
	"os"

	"github.com/donnie4w/go-logger/logger"
)

type Fileman struct {
	File         *os.File
	Fileinfo     os.FileInfo
	Blocksize    int64
	Blocknum     int64
	Lastpacksize int64
	Readoffset   int64
	Writeoffset  int64
}

var Default = Fileman{}

func Openfile(filename string) (f Fileman, err error) {

	f.File, err = os.Open(filename)
	if err != nil {
		logger.Error("create file err:", err)
		return
	}

	f.Fileinfo, err = f.File.Stat()
	f.Blocksize = 1024
	f.Blocknum = f.GetBlocknum()
	f.Lastpacksize = f.GetLastpacksize()
	f.Readoffset = 0
	f.Writeoffset = 0
	return f, err
}

func Newfile(filename string) (f Fileman, err error) {

	f.File, err = os.Create(filename)
	if err != nil {
		logger.Error("create file err:", err)
		return
	}

	f.Blocksize = 1024
	f.Readoffset = 0
	f.Writeoffset = 0
	return f, err
}

func (f *Fileman) GetBlocknum() int64 {
	return f.Getfilesize() / f.Blocksize
}
func (f *Fileman) GetLastpacksize() int64 {
	return f.Getfilesize() % f.Blocksize
}
func (f *Fileman) Getfilename() string {
	return f.Fileinfo.Name()
}

func (f *Fileman) Getfilesize() int64 {
	return f.Fileinfo.Size()
}

func (f *Fileman) Read(b []byte) (n int, err error) {
	n, err = f.File.ReadAt(b, f.Writeoffset)
	return
}

func (f *Fileman) Write(b []byte) (n int, err error) {
	n, err = f.File.WriteAt(b, f.Writeoffset)
	return
}
