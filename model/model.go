package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func PointInt64(v int64) *int64 {
	return &v
}

type Schedule struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Key  string             `bson:"key"`
	Jobs []Job              `bson:"jobs"`
}

type Job struct {
	Spec   string `bson:"spec"`
	Mode   string `bson:"mode"`
	Option bson.M `bson:"option"`
}

type HttpJob struct {
	Url     string                 `bson:"url"`
	Headers map[string]string      `bson:"headers"`
	Body    map[string]interface{} `bson:"body"`
}
