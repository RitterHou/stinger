package mylog

import (
	"log"
	"os"
)

var Info *log.Logger

func init() {
	logFile, err := os.OpenFile("local.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("打开日志文件失败：", err)
	}

	Info = log.New(logFile, "", log.LstdFlags|log.Lshortfile)
}
