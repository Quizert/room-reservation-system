package handlers

import (
	"context"
	"encoding/json"
	"github.com/Quizert/room-reservation-system/PaymentSystem/internal/models"
	"github.com/Quizert/room-reservation-system/PaymentSystem/internal/service"
	"log"
	"net/http"
	"time"
)

type PaymentHandler struct {
	PaymentService *service.PaymentService
}

func NewPaymentHandler(PaymentService *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		PaymentService: PaymentService,
	}
}

func (p *PaymentHandler) ProcessPayment(w http.ResponseWriter, r *http.Request) {
	var paymentRequest *models.PaymentRequest
	err := json.NewDecoder(r.Body).Decode(&paymentRequest)
	log.Println(paymentRequest)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		if err := p.PaymentService.ProcessPayment(ctx, paymentRequest); err != nil {
			log.Println("in handler payment processing failed:", err)
		}
	}()
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Payment processing started"))
}
