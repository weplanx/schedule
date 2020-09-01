package actions

import (
	"github.com/parnurzeal/gorequest"
	"schedule-microservice/app/types"
)

func Fetch(option types.FetchOption) (body []byte, errs []error) {
	agent := gorequest.New().Post(option.Url)
	if option.Headers != nil {
		for key, value := range option.Headers {
			agent.Set(key, value)
		}
	}
	if option.Body != nil {
		agent.Send(option.Body)
	}
	_, body, errs = agent.EndBytes()
	return
}
