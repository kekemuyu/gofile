package msg

import (
	"bytes"
	"encoding/binary"
	// "github.com/donnie4w/go-logger/logger"
)

type Msg struct {
	Id      uint32
	Datalen uint32
	Data    []byte
}

func (m *Msg) Pack(msg Msg) (b []byte, err error) {
	dataBuff := bytes.NewBuffer([]byte{})

	//写msgID
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.Id); err != nil {
		return nil, err
	}

	//写dataLen
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.Datalen); err != nil {
		return nil, err
	}

	//写data数据
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.Data); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

func (m *Msg) Unpack(b []byte) (msg Msg, err error) {

	//创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(b)

	var mg = Msg{}
	//读msgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return mg, err
	}

	//读dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Datalen); err != nil {
		return mg, err
	}

	//判断dataLen的长度是否超出我们允许的最大包长度

	//这里只需要把head的数据拆包出来就可以了，然后再通过head的长度，再从conn读取一次数据
	return msg, nil
}
