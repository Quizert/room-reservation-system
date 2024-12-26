package service

import (
	"context"
	"github.com/Quizert/room-reservation-system/AuthSvc/pkj/authpb"
	"github.com/Quizert/room-reservation-system/HotelSvc/api/grpc/hotelpb"
)

//go:generate mockgen -source=hotel.go -destination=mocks/hotel_mock.go -package=mocks
type HotelClient interface {
	GetRoomsByHotelId(ctx context.Context, req *hotelpb.GetRoomsRequest) (*hotelpb.GetRoomsResponse, error)
	GetOwnerIdByHotelId(ctx context.Context, req *hotelpb.GetOwnerIdRequest) (*hotelpb.GetOwnerIdResponse, error)
}

type AuthSvcClient interface {
	GetHotelierInformation(ctx context.Context, request *authpb.GetHotelierRequest) (*authpb.GetHotelierResponse, error)
}
