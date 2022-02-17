package api

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/weplanx/schedule/common"
)

type API struct {
	UnimplementedAPIServer
	*common.Inject
}

func (x *API) Create(ctx context.Context, req *CreateRequest) (_ *empty.Empty, err error) {
	return
}

func (x *API) Get(ctx context.Context, req *GetRequest) (rep *GetReply, err error) {
	return
}

func (x *API) Update(ctx context.Context, req *UpdateRequest) (_ *empty.Empty, err error) {
	return
}

func (x *API) Delete(ctx context.Context, req *DeleteRequest) (_ *empty.Empty, err error) {
	return
}
