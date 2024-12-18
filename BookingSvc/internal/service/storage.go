package service

import (
	"context"
	"github.com/Quizert/room-reservation-system/BookingSvc/internal/models"
	"time"
)

type Storage interface {
	CreateBooking(ctx context.Context, booking *models.Booking) (int, error)
	GetBookingsByUserID(ctx context.Context, userID int) ([]*models.Booking, error)
	GetBookingsByHotelID(ctx context.Context, bookingID int) (*models.Booking, error)

	UpdateBookingStatus(ctx context.Context, status string, bookingID int) error

	DeleteBooking(ctx context.Context, bookingID int) error

	GetUnavailableRoomsByHotelId(ctx context.Context, HotelID int, startDate, endDate time.Time) (map[int]struct{}, error)
}
