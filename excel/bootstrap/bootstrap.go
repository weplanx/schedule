package bootstrap

import (
	"github.com/caarlos0/env/v9"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/weplanx/workflow/excel/common"
	"net/http"
	"net/url"
)

func LoadValues() (values *common.Values, err error) {
	values = new(common.Values)
	if err = env.Parse(values); err != nil {
		return
	}
	return
}

func UseCos(values *common.Values) (client *cos.Client, err error) {
	u, _ := url.Parse(values.Cos.Url)
	b := &cos.BaseURL{BucketURL: u}
	client = cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  values.Cos.SecretId,
			SecretKey: values.Cos.SecretKey,
		},
	})
	return
}
