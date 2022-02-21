package api

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mitchellh/mapstructure"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/weplanx/schedule/app"
	"github.com/weplanx/schedule/common"
	"github.com/weplanx/schedule/model"
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

func (x *API) SetSchedule(key string, opts []model.Job) (err error) {
	jobs := make([]*app.Job, len(opts))
	for k, v := range opts {
		switch v.Mode {
		case "HTTP":
			var option model.HttpJob
			if err = mapstructure.Decode(v.Option, &option); err != nil {
				return
			}
			jobs[k] = app.HttpJob(v.Spec, option, x.Log)
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
			var count int64
			if count, err = x.Db.Collection(x.name()).CountDocuments(sessCtx, bson.M{
				"key": req.Key,
			}); err != nil {
				return
			}
			data := model.Schedule{Key: req.Key, Node: x.Values.Node}
			data.Jobs = make([]model.Job, len(req.Jobs))
			for k, v := range req.Jobs {
				var option map[string]interface{}
				if err = msgpack.Unmarshal(v.Option, &option); err != nil {
					return
				}
				data.Jobs[k] = model.Job{
					Spec:   v.Spec,
					Mode:   v.Mode,
					Option: option,
				}
			}
			if count != 0 {
				if _, err = x.Db.Collection(x.name(), wcMajorityCollectionOpts).
					ReplaceOne(sessCtx, bson.M{"key": req.Key}, data); err != nil {
					return
				}
				x.Schedule.Remove(req.Key)
			} else {
				if _, err = x.Db.Collection(x.name(), wcMajorityCollectionOpts).
					InsertOne(sessCtx, &data); err != nil {
					return
				}
			}
			if err = x.SetSchedule(req.Key, data.Jobs); err != nil {
				return
			}
			return
		},
	); err != nil {
		return
	}
	return &empty.Empty{}, nil
}

func (x *API) Get(ctx context.Context, req *GetRequest) (rep *GetReply, err error) {
	var cursor *mongo.Cursor
	if cursor, err = x.Db.Collection(x.name()).Find(ctx, bson.M{
		"key": bson.M{"$in": req.Keys},
	}); err != nil {
		return
	}
	var schedules []model.Schedule
	if err = cursor.All(ctx, &schedules); err != nil {
		return
	}
	rep = new(GetReply)
	rep.Data = make(map[string]*Schedule, len(schedules))
	for _, v := range schedules {
		key := v.Key
		state := x.Schedule.State(key)
		jobs := make([]*Job, len(v.Jobs))
		for kk, vv := range v.Jobs {
			var option []byte
			if option, err = msgpack.Marshal(vv.Option); err != nil {
				return
			}
			jobs[kk] = &Job{
				Spec:     vv.Spec,
				Mode:     vv.Mode,
				Option:   option,
				NextDate: model.PointInt64(state[kk].Next.Unix()),
				LastDate: model.PointInt64(state[kk].Prev.Unix()),
			}
		}
		rep.Data[key] = &Schedule{
			Key:  key,
			Jobs: jobs,
		}
	}
	return
}

func (x *API) Delete(ctx context.Context, req *DeleteRequest) (_ *empty.Empty, err error) {
	if _, err = x.Db.Collection(x.name()).DeleteOne(ctx, bson.M{
		"key": req.Key,
	}); err != nil {
		return
	}
	x.Schedule.Remove(req.Key)
	return &empty.Empty{}, nil
}
