package schema

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"schedule-microservice/app/types"
)

func (c *Schema) Update(config types.JobOption) (err error) {
	out, err := yaml.Marshal(config)
	if err != nil {
		return
	}
	if err = ioutil.WriteFile(
		c.autoload(config.Identity),
		out,
		0644,
	); err != nil {
		return
	}
	return
}
