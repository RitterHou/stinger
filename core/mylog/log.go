package mylog

import (
	"fmt"
	"github.com/ritterhou/stinger/core/common"
	"github.com/sirupsen/logrus"
	"github.com/t-tomalak/logrus-easy-formatter"
	"os"
	"runtime"
	"strings"
)

const evnFile = "env.txt"

var (
	separator = "\n"
	env       = "production"
)

func init() {
	if runtime.GOOS == "windows" {
		separator = "\r\n"
	}

	path := common.GetAbsPath(evnFile)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		content := common.ReadFile(path)
		env = string(content)
	}
}

func InitLog(file string, level string) {
	fmt.Printf("### Log file is %s, env is %s ###\n", file, env)

	if env == "develop" {
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors: false,
		})
		logrus.SetReportCaller(true)
	} else {
		logrus.SetFormatter(&easy.Formatter{
			TimestampFormat: "2006-01-02 15:04:05",
			LogFormat:       "%time% - [%lvl%] %msg%" + separator,
		})
	}

	logLevel := logrus.InfoLevel
	level = strings.ToUpper(level)
	switch level {
	case "DEBUG":
		logLevel = logrus.DebugLevel
	case "INFO":
		logLevel = logrus.InfoLevel
	case "WARN":
		logLevel = logrus.WarnLevel
	}
	logrus.SetLevel(logLevel)

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
