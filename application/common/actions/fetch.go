package actions

import (
	"github.com/parnurzeal/gorequest"
	"net/http"
	"time"
)

func Fetch(url string, headers map[string]string, body interface{}) ([]byte, []error) {
	agent := gorequest.New().Post(url)
	for key, value := range headers {
		agent.Set(key, value)
	}
	if body != nil {
		agent.Send(body)
	}
	_, resBody, errs := agent.
		Retry(3, 5*time.Second, http.StatusBadRequest, http.StatusInternalServerError).
		EndBytes()
	return resBody, errs
}
