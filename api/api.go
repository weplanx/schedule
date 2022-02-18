package api

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/weplanx/schedule/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"time"
)

type API struct {
	UnimplementedAPIServer
	*common.Inject
}

func (x *API) name() string {
	if x.Values.Database.Collection == "" {
		return "schedule"
	}
	return x.Values.Database.Collection
}

func (x *API) SetSchedule(key string, jobsOpt []map[string]interface{}) (err error) {
	var jobs []interface{}
	for _, v := range jobsOpt {
		switch v["mode"] {
		case "HTTP":
			jobs = append(jobs, common.HttpCallbackJob(
				v["spec"].(string),
				v["option"].(map[string]interface{}),
			))
			break
		}
		if err = x.Schedule.Set(key, jobs...); err != nil {
			return
		}
		x.Schedule.Start(key)
		return
	}
	return
}

func (x *API) Put(ctx context.Context, req *Schedule) (_ *empty.Empty, err error) {
	var session mongo.Session
	if session, err = x.Mongo.StartSession(); err != nil {
		return
	}
	defer session.EndSession(ctx)
	if _, err = session.WithTransaction(ctx,
		func(sessCtx mongo.SessionContext) (_ interface{}, err error) {
			wcMajority := writeconcern.New(writeconcern.WMajority(), writeconcern.WTimeout(time.Second))
			wcMajorityCollectionOpts := options.Collection().SetWriteConcern(wcMajority)
			var exists map[string]interface{}
			if err = x.Db.Collection(x.name()).FindOne(sessCtx, bson.M{
				"key": req.Key,
			}).Decode(&exists); err != nil {
				return
			}
			data := map[string]interface{}{"key": req.Key}
			jobsOpt := make([]map[string]interface{}, len(req.Jobs))
			for i, v := range req.Jobs {
				var option map[string]interface{}
				if err = msgpack.Unmarshal(v.Option, &option); err != nil {
					return
				}
				jobsOpt[i] = map[string]interface{}{
					"spec":   v.Spec,
					"mode":   v.Mode,
					"option": option,
				}
			}
			data["jobs"] = jobsOpt
			if len(exists) != 0 {
				if _, err = x.Db.Collection(x.name(), wcMajorityCollectionOpts).
					ReplaceOne(sessCtx, bson.M{"_id": exists["_id"]}, data); err != nil {
					return
				}
			} else {
				if _, err = x.Db.Collection(x.name(), wcMajorityCollectionOpts).
					InsertOne(sessCtx, data); err != nil {
					return
				}
			}
			x.Schedule.Remove(req.Key)
			if err = x.SetSchedule(req.Key, jobsOpt); err != nil {
				return
			}
			return
		},
	); err != nil {
		return
	}
	return
}

func (x *API) Get(ctx context.Context, req *GetRequest) (rep *GetReply, err error) {
	var cursor *mongo.Cursor
	if cursor, err = x.Db.Collection(x.name()).Find(ctx, bson.M{
		"key": bson.M{"$in": req.Keys},
	}); err != nil {
		return
	}
	var schedules []map[string]interface{}
	if err = cursor.All(ctx, &schedules); err != nil {
		return
	}
	rep = new(GetReply)
	data := make(map[string]*Schedule, len(schedules))
	for _, v := range schedules {
		key := v["key"].(string)
		jobs := make([]*Job, 0)
		state := x.Schedule.State(key)
		for ii, vv := range v["jobs"].([]map[string]interface{}) {
			jobs = append(jobs, &Job{
				Spec:     vv["spec"].(string),
				Mode:     vv["mode"].(string),
				Option:   vv["option"].([]byte),
				NextDate: common.PointInt64(state[ii].Next.Unix()),
				LastDate: common.PointInt64(state[ii].Prev.Unix()),
			})
		}
		data[key] = &Schedule{
			Key:  key,
			Jobs: jobs,
		}
	}
	rep.Data = data
	return
}

func (x *API) Delete(ctx context.Context, req *DeleteRequest) (_ *empty.Empty, err error) {
	if _, err = x.Db.Collection(x.name()).DeleteOne(ctx, bson.M{
		"key": req.Key,
	}); err != nil {
		return
	}
	x.Schedule.Remove(req.Key)
	return
}
