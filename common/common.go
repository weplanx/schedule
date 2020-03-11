package common

import "time"

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
		Identity string
		TimeZone string
		Start    bool
		Entries  map[string]*EntryOption
	}
	EntryOption struct {
		CronTime string
		Url      string
		Headers  map[string]string
		Body     interface{}
		NextDate time.Time
		LastDate time.Time
	}
)
