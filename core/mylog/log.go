package mylog

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/t-tomalak/logrus-easy-formatter"
	"os"
	"runtime"
)

var separator = "\n"

func init() {
	if runtime.GOOS == "windows" {
		separator = "\r\n"
	}
}

func InitLog(file string) {
	fmt.Printf("### Log file is %s ###\n", file)

	logrus.SetFormatter(&easy.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "%time% - [%lvl%] %msg%" + separator,
	})
	logrus.SetLevel(logrus.InfoLevel)

	if file == "stdout" {
		logrus.SetOutput(os.Stdout)
	} else {
		if _, err := os.Stat(file); !os.IsNotExist(err) {
			os.Remove(file)
		}
		logFile, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			fmt.Println("open log file failed", err)
		}
		logrus.SetOutput(logFile)
	}
}
