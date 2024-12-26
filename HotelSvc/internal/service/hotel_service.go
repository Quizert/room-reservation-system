package service

import (
	"errors"
	"github.com/Quizert/room-reservation-system/HotelSvc/internal/models"
)

type HotelRepository interface {
	GetAllHotels() ([]models.Hotel, error)
	AddHotel(hotel models.Hotel) error
	UpdateHotel(hotel models.Hotel) error
	GetHotelByID(id int) (*models.Hotel, error)
}

type HotelService struct {
	hotelRepo HotelRepository
}

// NewHotelService создает новый экземпляр HotelService.
func NewHotelService(hotelRepo HotelRepository) *HotelService {
	return &HotelService{hotelRepo: hotelRepo}
}

// GetAllHotels возвращает все отели в системе.
func (s *HotelService) GetAllHotels() ([]models.Hotel, error) {
	hotels, err := s.hotelRepo.GetAllHotels()
	if err != nil {
		return nil, err
	}
	return hotels, nil
}

// AddHotel добавляет новый отель.
func (s *HotelService) AddHotel(hotel models.Hotel) error {
	// Валидация данных отеля
	if hotel.Name == "" || hotel.OwnerId == 0 {
		return errors.New("invalid hotel data")
	}
	return s.hotelRepo.AddHotel(hotel)
}

// UpdateHotel обновляет информацию об отеле.
func (s *HotelService) UpdateHotel(hotel models.Hotel) error {
	// Проверка, существует ли отель
	existingHotel, err := s.hotelRepo.GetHotelByID(hotel.Id)
	if err != nil {
		return err
	}
	if existingHotel == nil {
		return errors.New("hotel not found")
	}

	return s.hotelRepo.UpdateHotel(hotel)
}

// GetHotelByID получает информацию об отеле по его ID.
func (s *HotelService) GetHotelByID(id int) (*models.Hotel, error) {
	return s.hotelRepo.GetHotelByID(id)
}
