package main

import (
	"fmt"
	"strconv"

	// "time"
	log "github.com/donnie4w/go-logger/logger"
)

type Tprocess struct {
	Filename string
	Filesize int64
	Runsize  int64
}

func Disprocessbar(rtdata Tprocess) {

	rate := (rtdata.Runsize * 100 / rtdata.Filesize)
	ratestr := strconv.Itoa(int(rate)) + "%"
	jsStr1 := fmt.Sprintf(`$('#curfilename').text("%s");$("#process").css("width","%s");$("#process").text("%s");`, rtdata.Filename, ratestr, ratestr)
	log.Debug(jsStr1)
	Defaultweb.UI.Eval(jsStr1)

}

// case <-t.C:
// 			size := (rtdata.Filesize - oldsize)
// 			if size <= 1024 {
// 				jsStr1 := fmt.Sprintf(`$('#speed').text(%s);`, size)
// 				Defaultweb.UI.Eval(jsStr1)
// 			} else if size <= 1024*1024 && size > 1024 {
// 				jsStr1 := fmt.Sprintf(`$('#speed').text(%f+"k");`, size/1024)
// 				Defaultweb.UI.Eval(jsStr1)
// 			} else {
// 				jsStr1 := fmt.Sprintf(`$('#speed').text(%f+"M");`, size/1024/1024)
// 				Defaultweb.UI.Eval(jsStr1)
// 			}
// 		}
