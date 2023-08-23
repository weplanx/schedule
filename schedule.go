package workflow

import (
	"errors"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/weplanx/workflow/typ"
	"time"
)

type Schedule struct {
	Namespace string
	Id        string
	Nats      *nats.Conn
	KeyValue  nats.KeyValue
}

func (x *Workflow) NewSchedule(id string) (schedule *Schedule, err error) {
	schedule = &Schedule{
		Namespace: x.Namespace,
		Id:        id,
		Nats:      x.Nats,
		KeyValue:  nil,
	}
	bucket := fmt.Sprintf(`%s_schedules_%s`, x.Namespace, id)
	if schedule.KeyValue, err = x.JetStream.KeyValue(bucket); err != nil {
		return
	}
	return
}

func (x *Schedule) Ping() (result bool, err error) {
	subj := fmt.Sprintf(`%s.schedules`, x.Namespace)
	var msg *nats.Msg
	if msg, err = x.Nats.Request(subj, []byte(x.Id), time.Second*5); err != nil {
		return
	}
	result = string(msg.Data) == "ok"
	return
}

func (x *Schedule) Lists() (keys []string, err error) {
	if keys, err = x.KeyValue.Keys(); err != nil {
		if errors.Is(err, nats.ErrNoKeysFound) {
			return []string{}, nil
		}
		return
	}
	return
}

func (x *Schedule) Get(key string) (option typ.ScheduleOption, err error) {
	var entry nats.KeyValueEntry
	if entry, err = x.KeyValue.Get(key); err != nil {
		return
	}
	if err = msgpack.Unmarshal(entry.Value(), &option); err != nil {
		return
	}
	var msg *nats.Msg
	subj := fmt.Sprintf(`%s.schedules.%s`, x.Namespace, x.Id)
	if msg, err = x.Nats.Request(subj, []byte(key), time.Second*3); err != nil {
		return
	}
	var states []typ.ScheduleState
	if err = msgpack.Unmarshal(msg.Data, &states); err != nil {
		return
	}
	for k, v := range states {
		option.Jobs[k].ScheduleState = v
	}
	return
}

func (x *Schedule) Set(key string, option typ.ScheduleOption) (err error) {
	var b []byte
	if b, err = msgpack.Marshal(option); err != nil {
		return
	}
	if _, err = x.KeyValue.Put(key, b); err != nil {
		return
	}
	return
}

func (x *Schedule) Status(key string, value bool) (err error) {
	var entry nats.KeyValueEntry
	if entry, err = x.KeyValue.Get(key); err != nil {
		return
	}
	var option typ.ScheduleOption
	if err = msgpack.Unmarshal(entry.Value(), &option); err != nil {
		return
	}
	option.Status = value
	return x.Set(key, option)
}

func (x *Schedule) Remove(key string) (err error) {
	return x.KeyValue.Delete(key)
}
