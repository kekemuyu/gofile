package main

import (
	"fmt"
	"strconv"

	"time"

	log "github.com/donnie4w/go-logger/logger"
)

type Tprocess struct {
	Filename string
	Filesize int64
	Runsize  int64
}

var btime time.Time
var oldsize int64

func Disprocessbar(rtdata Tprocess) {
	if rtdata.Runsize == 1024 {
		btime = time.Now()

	}

	rate := (rtdata.Runsize * 100 / rtdata.Filesize)
	ratestr := strconv.Itoa(int(rate)) + "%"

	jsStr1 := fmt.Sprintf(`$('#curfilename').text("%s");$("#process").css("width","%s");$("#process").text("%s");`, rtdata.Filename, ratestr, ratestr)

	log.Debug(jsStr1)
	Defaultweb.UI.Eval(jsStr1)
	if time.Since(btime) >= 1000000000 {
		btime = time.Now()

		size := (rtdata.Runsize - oldsize)
		log.Debug(size)
		if size < 1024 {
			sizestr := strconv.FormatInt(size, 10)
			jsStr1 := fmt.Sprintf(`$('#speed').text(%s);`, sizestr)
			Defaultweb.UI.Eval(jsStr1)
		} else if size < 1024*1024 && size >= 1024 {
			sizestr := strconv.FormatInt(size/1024, 10)
			jsStr1 := fmt.Sprintf(`$('#speed').text(%s+"k");`, sizestr)
			Defaultweb.UI.Eval(jsStr1)
		} else {
			sizestr := strconv.FormatInt(size/1024/1024, 10)
			jsStr1 := fmt.Sprintf(`$('#speed').text(%s+"M");`, sizestr)
			Defaultweb.UI.Eval(jsStr1)
		}
		oldsize = rtdata.Runsize
	}
}
