package log

import (
	plog "log"
)

func Debug(v ...interface{}) {
	return
	plog.Println(v)
}

func Error(v ...interface{}) {
	plog.Println(v)
}
