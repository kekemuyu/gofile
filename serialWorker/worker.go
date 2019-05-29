package serialworker

import (
	"bytes"
	"fmt"

	//	"fmt"
	"io"

	"github.com/jacobsa/go-serial/serial"
)

type Worker struct {
	Ioserial io.ReadWriteCloser
}

func New(port string) *Worker {
	opt := serial.OpenOptions{
		PortName:        port,
		BaudRate:        115200,
		DataBits:        8,
		StopBits:        1,
		ParityMode:      serial.PARITY_NONE,
		MinimumReadSize: 4,
	}

	s, err := serial.Open(opt)

	if err != nil {
		panic(err)
	}
	return &Worker{
		Ioserial: s,
	}
}

func (w *Worker) Read() bytes.Buffer {
	bs := make([]byte, 10240)
	count := 0
	var err error
	for {
		count, err = w.Ioserial.Read(bs)
		if err != nil {
			panic(err)
		}
		if count > 0 {
			fmt.Println(bs[:count])
			break
		}
	}

	var bb bytes.Buffer
	_, err = bb.Write(bs[:count])
	if err != nil {
		panic(err)
	}
	return bb
}

func (w *Worker) Write(bb bytes.Buffer) {
	fmt.Println(bb.Bytes())
	_, err := w.Ioserial.Write(bb.Bytes())
	if err != nil {
		panic(err)
	}
}
