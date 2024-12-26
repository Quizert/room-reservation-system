package grpc

import (
	"context"
	"fmt"
	"github.com/Quizert/room-reservation-system/AuthSvc/pkj/authpb"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type AuthSvcClient struct {
	Api  authpb.AuthServiceClient
	conn *grpc.ClientConn
}

func (a *AuthSvcClient) GetHotelierInformation(ctx context.Context, req *authpb.GetHotelierRequest) (*authpb.GetHotelierResponse, error) {
	return a.Api.GetHotelierInformation(ctx, req)
}

func NewAuthClient(grpcHost, grpcPort string) (*AuthSvcClient, error) {
	address := fmt.Sprintf("%s:%s", grpcHost, grpcPort)
	conn, err := grpc.Dial(
		address,
		grpc.WithInsecure(), // Рекомендуется использовать безопасное соединение (TLS) в продакшене
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		return nil, fmt.Errorf("could not connect: %w", err)
	}
	client := authpb.NewAuthServiceClient(conn)
	return &AuthSvcClient{Api: client, conn: conn}, nil
}

func (a *AuthSvcClient) Close() {
	err := a.conn.Close()
	if err != nil {
		fmt.Errorf("could not close connection: %w", err)
	}
}
