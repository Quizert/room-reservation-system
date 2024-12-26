package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Quizert/room-reservation-system/HotelSvc/internal/myerror"
)

type OwnerRepository interface {
	GetOwnerIdByHotelId(ctx context.Context, hotelID int) (int, error)
}

type OwnerService struct {
	ownerRepo OwnerRepository
}

func NewOwnerService(ownerRepo OwnerRepository) *OwnerService {
	return &OwnerService{ownerRepo: ownerRepo}
}

func (s *OwnerService) GetOwnerIdByHotelId(ctx context.Context, hotelID int) (int, error) {
	ownerID, err := s.ownerRepo.GetOwnerIdByHotelId(ctx, hotelID)
	if err != nil {
		if errors.Is(err, myerror.ErrHotelNotFound) {
			return 0, fmt.Errorf("in service GetOwnerIdByHotelId: %w", myerror.ErrHotelNotFound)
		}
		return 0, fmt.Errorf("in service GetOwnerIdByHotelId: %w", err)
	}
	return ownerID, err
}
