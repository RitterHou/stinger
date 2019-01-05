package common

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"testing"
)

var data = `
a: Easy!
b:
  c: 2
  d: [3, 4]
`

type Conf struct {
	A string
	B struct {
		RenamedC int   `yaml:"c"`
		D        []int `yaml:",flow"`
	}
}

func TestParse(t *testing.T) {
	c := Conf{}

	// 转化为结构体
	err := yaml.Unmarshal([]byte(data), &c)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- c:\n%v\n\n", c)

	// 转化为字节数组
	d, err := yaml.Marshal(&c)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- c dump:\n%s\n\n", string(d))

	m := make(map[interface{}]interface{})

	// 转化为map
	err = yaml.Unmarshal([]byte(data), &m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- m:\n%v\n\n", m)

	// 转回字节数组
	d, err = yaml.Marshal(&m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- m dump:\n%s\n\n", string(d))
}
