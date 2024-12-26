package service

import (
	"context"
	"github.com/Quizert/room-reservation-system/BookingSvc/internal/models"
	"time"
)

//go:generate mockgen -source=storage.go -destination=mocks/storage_mock.go -package=mocks
type Storage interface {
	CreateBooking(ctx context.Context, booking *models.BookingInfo) (int, error)
	GetBookingsByUserID(ctx context.Context, userID int) ([]*models.BookingInfo, error)
	GetBookingsByHotelID(ctx context.Context, hotelID int) ([]*models.BookingInfo, error)
	UpdateBookingStatus(ctx context.Context, status string, bookingID int) error

	GetUnavailableRoomsByHotelId(ctx context.Context, HotelID int, startDate, endDate time.Time) (map[int]struct{}, error)
}
