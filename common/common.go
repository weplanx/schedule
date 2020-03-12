package common

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type (
	AppOption struct {
		Debug  bool      `yaml:"debug"`
		Listen string    `yaml:"listen"`
		Log    LogOption `yaml:"log"`
	}
	LogOption struct {
		Storage    bool   `yaml:"storage"`
		StorageDir string `yaml:"storage_dir"`
		Socket     bool   `yaml:"socket"`
		SocketPort string `yaml:"socket_port"`
	}
	TaskOption struct {
		Identity string                  `yaml:"identity"`
		TimeZone string                  `yaml:"time_zone"`
		Start    bool                    `yaml:"start"`
		Entries  map[string]*EntryOption `yaml:"entries"`
	}
	EntryOption struct {
		CronTime string            `yaml:"cron_time"`
		Url      string            `yaml:"url"`
		Headers  map[string]string `yaml:"headers"`
		Body     interface{}       `yaml:"body"`
		NextDate time.Time         `yaml:"-"`
		LastDate time.Time         `yaml:"-"`
	}
)

func autoload(identity string) string {
	return "./config/autoload/" + identity + ".yml"
}

func ListConfig() (list []TaskOption, err error) {
	var files []os.FileInfo
	files, err = ioutil.ReadDir("./config/autoload")
	if err != nil {
		return
	}
	for _, file := range files {
		ext := filepath.Ext(file.Name())
		if ext == ".yml" {
			var in []byte
			in, err = ioutil.ReadFile("./config/autoload/" + file.Name())
			if err != nil {
				return
			}
			var config TaskOption
			err = yaml.Unmarshal(in, &config)
			if err != nil {
				return
			}
			list = append(list, config)
		}
	}
	return
}

func SaveConfig(data *TaskOption) (err error) {
	var out []byte
	out, err = yaml.Marshal(data)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(
		autoload(data.Identity),
		out,
		0644,
	)
	if err != nil {
		return
	}
	return
}

func RemoveConfig(identity string) error {
	return os.Remove(autoload(identity))
}
