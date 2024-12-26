package service

import (
	"context"
	"github.com/Quizert/room-reservation-system/BookingSvc/internal/models"
)

//go:generate mockgen -source=payment.go -destination=mocks/payment_mock.go -package=mocks
type PaymentSystemClient interface {
	CreatePaymentRequest(ctx context.Context, paymentRequest *models.PaymentRequest) error
}
