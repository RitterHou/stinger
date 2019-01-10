package resource

import (
	"github.com/rakyll/statik/fs"
	_ "github.com/ritterhou/stinger/local/statik"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"sync"
)

var (
	statikFS http.FileSystem
	contents map[string]string
	mutex    sync.Mutex
)

func init() {
	var err error
	statikFS, err = fs.New()
	if err != nil {
		logrus.Fatal(err)
	}
	contents = make(map[string]string)
}

// 根据文件名获取数据内容
func GetContent(filename string) string {
	if content, ok := contents[filename]; ok {
		return content
	} else {
		mutex.Lock()
		if content, ok = contents[filename]; ok {
			return content
		} else {
			file, err := statikFS.Open(filename)
			if err != nil {
				logrus.Warn(err)
				return ""
			}
			value, err := ioutil.ReadAll(file)
			if err != nil {
				logrus.Warn(err)
				return ""
			}
			content = string(value)
			contents[filename] = content
			return content
		}
		mutex.Unlock()
	}
	return ""
}
