package main

import (
	"context"
	"fmt"
	"github.com/Quizert/room-reservation-system/BookingSvc/internal/app"
)

func main() {
	service := app.NewApp()
	ctxWithCancel, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := service.Init(ctxWithCancel)
	if err != nil {
		return
	}

	err = service.Start(ctxWithCancel)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
