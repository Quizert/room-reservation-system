package models

type PaymentRequest struct {
	CardNumber string `json:"card_number"`
	Amount     int    `json:"amount"`
	WebHookURL string `json:"web_hook_url"`

	MetaData map[string]interface{} `json:"meta_data"` // Произвольные метаданные
}

type PaymentResponse struct {
	Status string `json:"status"`

	MetaData map[string]interface{} `json:"meta_data"` // Произвольные метаданные
}
