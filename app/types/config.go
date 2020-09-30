package types

type Config struct {
	Debug   string        `yaml:"debug"`
	Listen  string        `yaml:"listen"`
	Logging LoggingOption `yaml:"logging"`
}
