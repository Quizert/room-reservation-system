package grpc

import (
	"context"
	"fmt"
	"github.com/Quizert/room-reservation-system/HotelSvc/api/grpc/hotelpb"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type HotelSvcClient struct {
	Api  hotelpb.HotelServiceClient
	conn *grpc.ClientConn
}

func (c *HotelSvcClient) GetOwnerIdByHotelId(ctx context.Context, req *hotelpb.GetOwnerIdRequest) (*hotelpb.GetOwnerIdResponse, error) {
	return c.Api.GetOwnerIdByHotelId(ctx, req)
}

func (c *HotelSvcClient) GetRoomsByHotelId(ctx context.Context, req *hotelpb.GetRoomsRequest) (*hotelpb.GetRoomsResponse, error) {
	return c.Api.GetRoomsByHotelId(ctx, req)
}

func NewHotelClient(grpcHost, grpcPort string) (*HotelSvcClient, error) {
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", grpcHost, grpcPort),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		return nil, fmt.Errorf("could not connect: %w", err)
	}
	client := hotelpb.NewHotelServiceClient(conn)
	return &HotelSvcClient{Api: client, conn: conn}, nil
}

func (c *HotelSvcClient) Close() {
	err := c.conn.Close()
	if err != nil {
		fmt.Errorf("could not close connection: %w", err)
	}
}
