package mylog

import (
	"fmt"
	"log"
	"os"
)

func InitLog(file string) {
	fmt.Printf("### Log file is %s ###\n", file)

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if file == "stdout" {
		log.SetOutput(os.Stdout)
		return
	}

	if _, err := os.Stat(file); !os.IsNotExist(err) {
		os.Remove(file)
	}
	logFile, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalln("open log file failed", err)
	}
	log.SetOutput(logFile)
}
