package internal

import (
	"bytes"
	"encoding/json"

	"os"
)

type File struct {
	Name string
	Size int64
	Data []byte
}

type Control interface {
	Reader() bytes.Buffer
	Writer(bytes.Buffer)
}

type Buffer struct{}

var Defaultbuffer Buffer

func (b Buffer) GetBytesbuffer(fileName string) bytes.Buffer {

	fin, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer fin.Close()

	fileInfo, err := fin.Stat()
	if err != nil {
		panic(err)
	}

	var bs []byte

	_, err = fin.Read(bs)
	if err != nil {
		panic(err)
	}

	fileObj := File{
		Name: fileInfo.Name(),
		Size: fileInfo.Size(),
		Data: bs,
	}
	objBytes, err := json.Marshal(fileObj)
	if err != nil {
		panic(err)
	}

	var data bytes.Buffer
	_, err = data.Write(objBytes)
	if err != nil {
		panic(err)
	}
}
