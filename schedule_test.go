package workflow_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/weplanx/workflow/typ"
	"testing"
)

func TestSchedule_Set(t *testing.T) {
	err := schedule.Set("api", typ.ScheduleOption{
		Status: false,
		Jobs: []typ.ScheduleJob{
			{
				Mode: "HTTP",
				Spec: "*/2 * * * * *",
				Option: typ.HttpOption{
					Url: "https://api.kainonly.com/whoami",
				},
			},
		},
	})
	assert.NoError(t, err)
}

func TestSchedule_List(t *testing.T) {
	keys, err := schedule.List()
	assert.NoError(t, err)
	t.Log(keys)
}

func TestSchedule_Get(t *testing.T) {
	jobs, err := schedule.Get("api")
	assert.NoError(t, err)
	t.Log(jobs)
}

func TestSchedule_StatusStart(t *testing.T) {
	err := schedule.Status("api", true)
	assert.NoError(t, err)
}

func TestSchedule_StatusStop(t *testing.T) {
	err := schedule.Status("api", false)
	assert.NoError(t, err)
}

func TestSchedule_Remove(t *testing.T) {
	err := schedule.Remove("api")
	assert.NoError(t, err)
}
