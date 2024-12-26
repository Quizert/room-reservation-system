package service

import "github.com/Quizert/room-reservation-system/HotelSvc/api/grpc/hotelpb"

type RoomRepository interface {
	GetRoomsByHotelId(id int) ([]*hotelpb.Room, error)
}

type RoomService struct {
	roomRepo RoomRepository
}

// NewRoomService создает новый экземпляр RoomService.
func NewRoomService(roomRepo RoomRepository) *RoomService {
	return &RoomService{roomRepo: roomRepo}
}

func (s *RoomService) GetRoomsByHotelId(id int) ([]*hotelpb.Room, error) {
	rooms, err := s.roomRepo.GetRoomsByHotelId(id)
	if err != nil {
		return nil, err
	}
	return rooms, err
}
