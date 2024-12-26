package main

import (
	"context"
	"fmt"
	"github.com/Quizert/room-reservation-system/AuthSvc/internal/app"
)

func main() {
	app := app.NewApp()
	ctxWithCancel, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := app.Init(ctxWithCancel)
	if err != nil {
		return
	}

	err = app.Start(ctxWithCancel)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
