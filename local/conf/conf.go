package conf

import (
	"github.com/ritterhou/stinger/core/common"
	"log"
)

var conf = make(map[interface{}]interface{})

func LoadConf(filename string) {
	path := common.GetAbsPath(filename)
	content := common.ReadFile(path)
	conf = common.MarshalYaml(content)
	log.Println("Load and marshal local configuration.")
}

func GetConf() map[interface{}]interface{} {
	return conf
}
