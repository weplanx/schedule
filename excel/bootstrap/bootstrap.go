package bootstrap

import (
	"excel/common"
	"github.com/caarlos0/env/v9"
	"github.com/tencentyun/cos-go-sdk-v5"
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
	option := values.Cos
	var u *url.URL
	u, err = url.Parse(option.Url)
	b := &cos.BaseURL{BucketURL: u}
	client = cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  option.SecretId,
			SecretKey: option.SecretKey,
		},
	})
	return
}
