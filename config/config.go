package config

import (
	"schedule-microservice/config/options"
)

type Config struct {
	Debug    string                 `yaml:"debug"`
	Listen   string                 `yaml:"listen"`
	Gateway  string                 `yaml:"gateway"`
	Transfer options.TransferOption `yaml:"transfer"`
	Filelog  string                 `yaml:"filelog"`
}
