package schema

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
	"schedule-microservice/app/types"
)

func (c *Schema) Lists() (options []types.JobOption, err error) {
	var files []os.FileInfo
	files, err = ioutil.ReadDir(c.path)
	if err != nil {
		return
	}
	for _, file := range files {
		ext := filepath.Ext(file.Name())
		if ext == ".yml" {
			var bytes []byte
			bytes, err = ioutil.ReadFile(c.path + file.Name())
			if err != nil {
				return
			}
			var option types.JobOption
			err = yaml.Unmarshal(bytes, &option)
			if err != nil {
				return
			}
			options = append(options, option)
		}
	}
	return
}
