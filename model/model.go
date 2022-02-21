package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func PointInt64(v int64) *int64 {
	return &v
}

type Schedule struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Key  string             `bson:"key"`
	Node string             `bson:"node"`
	Jobs []Job              `bson:"jobs"`
}

type Job struct {
	Spec   string                 `bson:"spec"`
	Mode   string                 `bson:"mode"`
	Option map[string]interface{} `bson:"option"`
}

type HttpJob struct {
	Url     string                 `bson:"url"`
	Headers map[string]string      `bson:"headers"`
	Body    map[string]interface{} `bson:"body"`
}
