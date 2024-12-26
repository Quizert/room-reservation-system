package controller

import (
	"github.com/Quizert/room-reservation-system/PaymentSystem/internal/controller/handlers"
	"net/http"
)

func SetupRoutes(PaymentHandler *handlers.PaymentHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/payment", PaymentHandler.ProcessPayment)
	return mux
}
