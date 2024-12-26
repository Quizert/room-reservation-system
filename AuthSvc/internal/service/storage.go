package service

import (
	"context"
	"github.com/Quizert/room-reservation-system/AuthSvc/internal/models"
	"github.com/Quizert/room-reservation-system/AuthSvc/pkj/authpb"
)

type Storage interface {
	RegisterUser(ctx context.Context, user *models.User) (int, error)
	LoginUser(ctx context.Context, chatID string) (*models.User, error)
	IsHotelier(ctx context.Context, userID int) (bool, error)
	GetHotelierInformation(ctx context.Context, request *authpb.GetHotelierRequest) (*authpb.GetHotelierResponse, error)
}
