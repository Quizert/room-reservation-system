package models

import "strconv"

type PaymentRequest struct {
	CardNumber string `json:"card_number"`
	Amount     int    `json:"amount"`
	WebHookURL string `json:"web_hook_url"`

	MetaData *BookingMessage `json:"meta_data"` //Это в meta data
}

type PaymentResponse struct {
	Status string `json:"status"`

	MetaData *BookingMessage `json:"meta_data"` //Это в meta data
}

func ToPaymentRequest(bookingMessage *BookingMessage, cardNumber string, amount int) *PaymentRequest {
	return &PaymentRequest{
		CardNumber: cardNumber,
		Amount:     amount,
		WebHookURL: "http://booking-service:8080/bookings/payment/response?booking_id=" + strconv.Itoa(bookingMessage.BookingID),

		MetaData: bookingMessage,
	}
}
