package manage

import (
	"schedule-microservice/app/utils"
)

type JobsManager struct {
	options *utils.SyncMapJobOption
	runtime *utils.SyncMapCron
	entries *utils.SyncMapEntryID
}

func NewJobsManager() (manager *JobsManager, err error) {
	c := new(JobsManager)
	c.options = utils.NewSyncMapJobOption()
	c.runtime = utils.NewSyncMapCron()
	c.entries = utils.NewSyncMapEntryID()
	return
}
