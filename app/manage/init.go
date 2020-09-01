package manage

import (
	"errors"
	"github.com/robfig/cron/v3"
	"schedule-microservice/app/types"
	"schedule-microservice/app/utils"
)

type JobsManager struct {
	options    *utils.SyncMapJobOption
	runtime    *utils.SyncMapCron
	entryIDSet *utils.SyncMapEntryID
}

func NewJobsManager() (manager *JobsManager, err error) {
	c := new(JobsManager)
	c.options = utils.NewSyncMapJobOption()
	c.runtime = utils.NewSyncMapCron()
	c.entryIDSet = utils.NewSyncMapEntryID()
	return
}

func (c *JobsManager) empty(identity string) error {
	if c.options.Empty(identity) || c.runtime.Empty(identity) {
		return errors.New("this identity does not exists")
	}
	return nil
}

func (c *JobsManager) GetIdentityCollection() []string {
	var keys []string
	for key := range c.options.Map {
		keys = append(keys, key)
	}
	return keys
}

func (c *JobsManager) GetOption(identity string) (option *types.JobOption, err error) {
	if err = c.empty(identity); err != nil {
		return
	}
	option = c.options.Get(identity)
	return
}

func (c *JobsManager) GetJob(identity string) (job *cron.Cron, err error) {
	if err = c.empty(identity); err != nil {
		return
	}
	job = c.runtime.Get(identity)
	return
}

func (c *JobsManager) GetEntryIDSet(identity string) (entryIDs map[string]cron.EntryID, err error) {
	if err = c.empty(identity); err != nil {
		return
	}
	entryIDs = c.entryIDSet.Map[identity]
	return
}

func (c *JobsManager) termination(identity string) {
	if c.runtime.Empty(identity) || c.options.Empty(identity) {
		return
	}
	runtime := c.runtime.Get(identity)
	for _, entry := range runtime.Entries() {
		runtime.Remove(entry.ID)
	}
	runtime.Stop()
}
