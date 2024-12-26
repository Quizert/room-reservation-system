package service

import "context"

//go:generate mockgen -source=producer.go -destination=mocks/producer_mock.go -package=mocks
type MessageProducer interface {
	SendUserMessage(ctx context.Context, value []byte) error
	SendHotelierMessage(ctx context.Context, value []byte) error
}
