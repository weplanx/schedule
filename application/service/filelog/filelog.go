package filelog

import (
	"errors"
	"go.uber.org/zap"
	"os"
	"path"
	"time"
)

var (
	NotSet = errors.New("please set the filelog path")
)

type Filelog struct {
	on   bool
	path string
}

func New(path string) *Filelog {
	c := new(Filelog)
	c.on = path != ""
	c.path = path
	if _, err := os.Stat(c.path); os.IsNotExist(err) {
		os.Mkdir(c.path, os.ModeDir)
	}
	return c
}

func (c *Filelog) NewLogger(identity string) (logger *zap.Logger, err error) {
	if !c.on {
		return nil, NotSet
	}
	cfg := zap.NewProductionConfig()
	outputId := path.Join(c.path, identity)
	if _, err = os.Stat(outputId); os.IsNotExist(err) {
		os.Mkdir(outputId, os.ModeDir)
	}
	now := time.Now().Format("2006-01-02")
	cfg.OutputPaths = []string{path.Join(outputId, now+".log")}
	return cfg.Build()
}
