package config

import (
	"os"
)

type Config struct {
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	GRPCHotelHost    string
	GRPCHotelPort    string
	GRPCAuthHost     string
	GRPCAuthPort     string
	HTTPPort         string
	HTTPMetricPort   string
	KafkaBroker      string
	KafkaTopicClient string
	KafkaTopicHotel  string
	PaymentSvcURL    string
}

func LoadConfig() (*Config, error) {
	return &Config{
		DBHost:           os.Getenv("BOOKING_DB_HOST"),
		DBPort:           os.Getenv("BOOKING_DB_PORT"),
		DBUser:           os.Getenv("BOOKING_DB_USER"),
		DBPassword:       os.Getenv("BOOKING_DB_PASSWORD"),
		DBName:           os.Getenv("BOOKING_DB_NAME"),
		GRPCHotelHost:    os.Getenv("HOTEL_GRPC_HOST"),
		GRPCHotelPort:    os.Getenv("HOTEL_GRPC_PORT"),
		GRPCAuthHost:     os.Getenv("AUTH_GRPC_HOST"),
		GRPCAuthPort:     os.Getenv("AUTH_GRPC_PORT"),
		HTTPPort:         os.Getenv("BOOKING_HTTP_PORT"),
		HTTPMetricPort:   os.Getenv("BOOKING_HTTP_METRIC_PORT"),
		KafkaBroker:      os.Getenv("KAFKA_BROKER"),
		KafkaTopicClient: os.Getenv("KAFKA_TOPIC_CLIENT"),
		KafkaTopicHotel:  os.Getenv("KAFKA_TOPIC_HOTEL"),
		PaymentSvcURL:    os.Getenv("PAYMENT_SERVICE_URL"),
	}, nil
}
