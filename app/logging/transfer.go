package logging

import (
	"context"
	"encoding/json"
	"errors"
	"schedule-microservice/app/types"
	pb "schedule-microservice/transfer"
)

func (c *Logging) forTransfer(push *types.LoggingPush) (err error) {
	var data []byte
	data, err = json.Marshal(push.Message)
	if err != nil {
		return
	}
	response, err := c.transfer.Push(context.Background(), &pb.PushParameter{
		Identity: c.Transfer.Id,
		Data:     data,
	})
	if err != nil {
		return
	}
	if response.Error != 0 {
		return errors.New(response.Msg)
	}
	return
}
