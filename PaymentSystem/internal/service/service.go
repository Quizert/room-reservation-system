package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Quizert/room-reservation-system/PaymentSystem/internal/models"
	"log"
	"net/http"
	"time"
)

type PaymentService struct {
	client *http.Client
}

func NewPaymentService() *PaymentService {
	return &PaymentService{client: &http.Client{}}
}

func (p *PaymentService) sendWebHook(ctx context.Context, url string, paymentResponse *models.PaymentResponse) error {
	data, err := json.Marshal(paymentResponse)
	if err != nil {
		return fmt.Errorf("failed to marshal paymentResponse: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	log.Println(string(data), "ADASDASDKLASFKOL:ASJFOLIKASFJHIKOPASFJKLASJKO")
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	return nil
}

func (p *PaymentService) ProcessPayment(ctx context.Context, req *models.PaymentRequest) error {
	time.Sleep(15 * time.Second) // Имитация обратки платежа, связь с банком и т.д.
	paymentResponse := &models.PaymentResponse{
		Status: "success",

		MetaData: req.MetaData,
	}

	select {
	case <-ctx.Done():
		paymentResponse.Status = "failed"
		сtx := context.Background()
		err := p.sendWebHook(сtx, req.WebHookURL, paymentResponse)
		if err != nil {
			return fmt.Errorf("payment process myerror: %w", err)
		}
	default:
		err := p.sendWebHook(ctx, req.WebHookURL, paymentResponse)
		if err != nil {
			return fmt.Errorf("payment process myerror: %w", err)
		}
	}
	return nil
}
