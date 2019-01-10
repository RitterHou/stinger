package conf

import (
	"github.com/ritterhou/stinger/core/common"
)

var conf = make(map[interface{}]interface{})

func LoadConf(filename string) {
	path := common.GetAbsPath(filename)
	content := common.ReadFile(path)
	conf = common.MarshalYaml(content)
}

func GetConf() map[interface{}]interface{} {
	return conf
}
