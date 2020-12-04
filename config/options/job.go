package options

type JobOption struct {
	Identity string                  `yaml:"identity"`
	TimeZone string                  `yaml:"time_zone"`
	Start    bool                    `yaml:"start"`
	Entries  map[string]*EntryOption `yaml:"entries"`
}
