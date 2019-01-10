package common

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
)

var (
	root string
)

func init() {
	dir, err := os.Getwd()
	if err != nil {
		logrus.Fatal(err)
	}
	root = dir
}

func GetAbsPath(filename string) (absPath string) {
	absPath = path.Join(root, filename)
	return
}

func ReadFile(filePath string) []byte {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		logrus.Fatal(err)
	}
	return data
}

func WriteFile(filePath string, data []byte) {
	ioutil.WriteFile(filePath, data, 0644)
}
