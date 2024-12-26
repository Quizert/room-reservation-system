package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

type Producer struct {
	writer        *kafka.Writer
	userTopic     string
	hotelierTopic string
}

func NewProducer(brokers []string, userTopic, hotelierTopic string) *Producer {
	return &Producer{
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:      brokers,
			Balancer:     &kafka.LeastBytes{},
			BatchTimeout: 10 * time.Microsecond,
		}),
		userTopic:     userTopic,
		hotelierTopic: hotelierTopic,
	}
}

func (p *Producer) sendMessage(ctx context.Context, topic string, value []byte) error {
	msg := kafka.Message{
		Topic: topic,
		Value: value,
	}
	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		log.Printf("Failed to send message to topic %s: %v", topic, err)
		return err
	}
	log.Printf("Message sent to topic %s: %s", topic, value)
	return nil
}

func (p *Producer) SendUserMessage(ctx context.Context, value []byte) error {
	return p.sendMessage(ctx, p.userTopic, value)
}

func (p *Producer) SendHotelierMessage(ctx context.Context, value []byte) error {
	return p.sendMessage(ctx, p.hotelierTopic, value)
}
func (p *Producer) Close() error {
	return p.writer.Close()
}
