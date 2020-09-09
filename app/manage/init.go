package manage

import (
	"errors"
	"github.com/robfig/cron/v3"
	"schedule-microservice/app/logging"
	"schedule-microservice/app/schema"
	"schedule-microservice/app/types"
	"schedule-microservice/app/utils"
)

type JobsManager struct {
	options    *utils.SyncMapJobOption
	runtime    *utils.SyncMapCron
	entryIDSet *utils.SyncMapEntryID
	logging    *logging.Logging
	schema     *schema.Schema
}

func NewJobsManager(schema *schema.Schema, logging *logging.Logging) (manager *JobsManager, err error) {
	manager = new(JobsManager)
	manager.options = utils.NewSyncMapJobOption()
	manager.runtime = utils.NewSyncMapCron()
	manager.entryIDSet = utils.NewSyncMapEntryID()
	manager.logging = logging
	manager.schema = schema
	var jobOptions []types.JobOption
	jobOptions, err = manager.schema.Lists()
	if err != nil {
		return
	}
	for _, option := range jobOptions {
		err = manager.Put(option)
		if err != nil {
			return
		}
	}
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
