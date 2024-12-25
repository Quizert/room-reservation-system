package config

import (
	"os"
)

type Config struct {
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	GRPCHost       string
	GRPCPort       string
	HTTPPort       string
	HTTPMetricPort string
	TokenTTl       string
	Secret         string
}

func LoadConfig() (*Config, error) {
	return &Config{
		DBHost:         os.Getenv("AUTH_DB_HOST"),
		DBPort:         os.Getenv("AUTH_DB_PORT"),
		DBUser:         os.Getenv("AUTH_DB_USER"),
		DBPassword:     os.Getenv("AUTH_DB_PASSWORD"),
		DBName:         os.Getenv("AUTH_DB_NAME"),
		GRPCHost:       os.Getenv("AUTH_GRPC_HOST"),
		GRPCPort:       os.Getenv("AUTH_GRPC_PORT"),
		HTTPPort:       os.Getenv("AUTH_HTTP_PORT"),
		HTTPMetricPort: os.Getenv("AUTH_HTTP_METRIC_PORT"),
		Secret:         os.Getenv("AUTH_SECRET"),
		TokenTTl:       os.Getenv("AUTH_TOKEN_TTL"),
	}, nil
}
