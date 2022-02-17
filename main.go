package main

import (
	"github.com/weplanx/schedule/bootstrap"
	"net"
)

func main() {
	v, err := bootstrap.SetValues()
	if err != nil {
		panic(err)
	}
	lis, err := net.Listen("tcp", v.Address)
	if err != nil {
		panic(err)
	}
	app, err := App(v)
	if err != nil {
		panic(err)
	}
	app.Serve(lis)
}
