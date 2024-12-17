package grpc

import (
	"context"
	"fmt"
	"github.com/Quizert/room-reservation-system/HotelSvc/api/grpc/hotelpb"
	"google.golang.org/grpc"
)

type HotelSvcClient struct {
	Api  hotelpb.HotelServiceClient
	conn *grpc.ClientConn
}

func (c *HotelSvcClient) GetRoomsByHotelId(ctx context.Context, req *hotelpb.GetRoomsRequest) (*hotelpb.GetRoomsResponse, error) {
	return c.Api.GetRoomsByHotelId(ctx, req)
}

func NewHotelClient(grpcHost, grpcPort string) (*HotelSvcClient, error) {
	address := fmt.Sprintf("%s:%s", grpcHost, grpcPort)
	conn, err := grpc.Dial(address, grpc.WithInsecure()) // Добавить ретраи мб сервис упадет??
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