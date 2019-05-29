package pipe

import (
	"bytes"
)

type Control interface {
	Read() bytes.Buffer
	Write(bb bytes.Buffer)
}
