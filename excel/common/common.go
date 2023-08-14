package common

import "github.com/tencentyun/cos-go-sdk-v5"

type Inject struct {
	Values *Values
	Client *cos.Client
}

type Values struct {
	Address string `env:"ADDRESS" envDefault:":9000"`
	Cos     struct {
		Url       string `env:"URL"`
		SecretId  string `env:"SECRETID"`
		SecretKey string `env:"SECRETKEY"`
	} `envPrefix:"COS_"`
}
