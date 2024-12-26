package paymentsvc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Quizert/room-reservation-system/BookingSvc/internal/models"
	"log"
	"net/http"
	"time"
)

type Client struct {
	baseUrl string
	client  *http.Client
}

func NewPaymentSvcClient(baseUrl string) *Client {
	return &Client{
		baseUrl: baseUrl,
		client:  &http.Client{Timeout: 5 * time.Minute},
	}
}

func (c *Client) CreatePaymentRequest(ctx context.Context, paymentRequest *models.PaymentRequest) error {
	jsonRequest, err := json.Marshal(paymentRequest)
	if err != nil {
		return fmt.Errorf("myerror in marshaling json: %w", err)
	}
	log.Println("JSON PaymentRequest: ", string(jsonRequest))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseUrl, bytes.NewBuffer(jsonRequest))
	if err != nil {
		return fmt.Errorf("myerror in creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("myerror in sending request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("myerror in payment service status: %s", resp.Status)
	}
	return nil
}
