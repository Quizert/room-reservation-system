package main

import (
	"github.com/Quizert/room-reservation-system/PaymentSystem/internal/controller"
	"github.com/Quizert/room-reservation-system/PaymentSystem/internal/controller/handlers"
	"github.com/Quizert/room-reservation-system/PaymentSystem/internal/service"
	"log"
	"net/http"
)

type App struct {
	service *service.PaymentService
	server  *http.Server
}

func NewApp() *App {
	return &App{}
}

func main() {
	app := NewApp()
	app.service = service.NewPaymentService()

	paymentHandler := handlers.NewPaymentHandler(app.service)
	app.server = &http.Server{
		Addr:    ":8080",
		Handler: controller.SetupRoutes(paymentHandler),
	}

	err := app.server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
