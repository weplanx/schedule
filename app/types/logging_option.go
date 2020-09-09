package types

type LoggingOption struct {
	Storage  string         `yaml:"storage"`
	Transfer TransferOption `yaml:"transfer"`
}
