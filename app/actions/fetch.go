package actions

import (
	"github.com/parnurzeal/gorequest"
	"net/http"
	"schedule-microservice/app/types"
	"time"
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
	_, body, errs = agent.
		Retry(3, 5*time.Second, http.StatusBadRequest, http.StatusInternalServerError).
		EndBytes()
	return
}
