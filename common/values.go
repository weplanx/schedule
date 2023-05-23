package common

import (
	"errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

func SetValues() (values *Values, err error) {
	if _, err = os.Stat("./config/config.yml"); os.IsNotExist(err) {
		err = errors.New("静态配置不存在，请检查路径 [./config/config.yml]")
		return
	}
	var b []byte
	b, err = ioutil.ReadFile("./config/config.yml")
	if err != nil {
		return
	}
	err = yaml.Unmarshal(b, &values)
	if err != nil {
		return
	}
	return
}

type Values struct {
	Debug     bool   `yaml:"debug"`
	Namespace string `yaml:"namespace"`
	Nats      struct {
		Hosts []string `yaml:"hosts"`
		Nkey  string   `yaml:"nkey"`
	} `yaml:"nats"`
}
