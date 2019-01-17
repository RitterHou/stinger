package conf

import (
	"github.com/ritterhou/stinger/core/common"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Conf struct {
	ServerPort int    `yaml:"server_port"`
	Password   string `yaml:"password"`
	LogFile    string `yaml:"log_file"`
	LogLevel   string `yaml:"log_level"`
}

var conf Conf

func LoadConf(filename string) {
	path := common.GetAbsPath(filename)
	content := common.ReadFile(path)

	err := yaml.Unmarshal(content, &conf)
	if err != nil {
		logrus.Fatal(err)
	}

	if conf.ServerPort == 0 {
		conf.ServerPort = 26800
	}
	if conf.Password == "" {
		conf.Password = "123456"
	}
	if conf.LogFile == "" {
		conf.LogFile = "stdout"
	}
	if conf.LogLevel == "" {
		conf.LogLevel = "WARN"
	}
}

func GetConf() Conf {
	return conf
}
