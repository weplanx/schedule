package main

import (
	"context"
	"github.com/weplanx/workflow/worker/bootstrap"
	"os"
	"os/signal"
	"time"
)

func main() {
	app, err := bootstrap.NewApp()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err = app.Run(ctx); err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
