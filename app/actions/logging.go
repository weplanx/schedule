package actions

import (
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"schedule-microservice/app/types"
	"time"
)

func Logging(option *types.LoggingOption, push *types.LoggingPush) (err error) {
	if option.Storage != "" {
		err = loggingForStorage(option, push)
		if err != nil {
			return
		}
	}
	return
}

func loggingForStorage(option *types.LoggingOption, push *types.LoggingPush) (err error) {
	logger := logrus.New()
	if _, err := os.Stat(option.Storage); os.IsNotExist(err) {
		os.Mkdir(option.Storage, os.ModeDir)
	}
	identityPath := path.Join(option.Storage, push.Identity)
	if _, err := os.Stat(identityPath); os.IsNotExist(err) {
		os.Mkdir(identityPath, os.ModeDir)
	}
	date := time.Now().Format("2006-01-02")
	filename := path.Join(option.Storage, push.Identity, date+".log")
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
