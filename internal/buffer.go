package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"os"
)

type File struct {
	Name string
	Size int64
	Data []byte
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
	bs, err = ioutil.ReadFile(fileName)

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
	return data
}

func (b Buffer) PutBytesbufferToFile(bs []byte) {
	var obj File
	err := json.Unmarshal(bs, &obj)
	if err != nil {
		fmt.Println(err)
		return
	}
	// f, err := os.Create(obj.Name)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// defer f.Close()
	ioutil.WriteFile(obj.Name, obj.Data, 0666)
	// _, err = f.Write(obj.Data)
	// if err != nil {
	// 	panic(err)
	// }

}
