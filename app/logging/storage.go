package logging

import (
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"schedule-microservice/app/types"
	"time"
)

func (c *Logging) forStorage(push *types.LoggingPush) (err error) {
	logger := logrus.New()
	if _, err := os.Stat(c.Storage); os.IsNotExist(err) {
		os.Mkdir(c.Storage, os.ModeDir)
	}
	identityPath := path.Join(c.Storage, push.Identity)
	if _, err := os.Stat(identityPath); os.IsNotExist(err) {
		os.Mkdir(identityPath, os.ModeDir)
	}
	date := time.Now().Format("2006-01-02")
	filename := path.Join(c.Storage, push.Identity, date+".log")
	var file *os.File
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		file, err = os.Create(filename)
		if err != nil {
			return
		}
	} else {
		file, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			return
		}
	}
	logger.SetOutput(file)
	if !push.HasError {
		logger.Info(push.Message)
	} else {
		logger.Error(push.Message)
	}
	return
}
