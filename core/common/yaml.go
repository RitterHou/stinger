package common

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// 将yaml格式的字符串转化为map类型
func MarshalYaml(source []byte) map[interface{}]interface{} {
	m := make(map[interface{}]interface{})
	err := yaml.Unmarshal(source, &m)
	if err != nil {
		logrus.Fatal(err)
	}
	return m
}
