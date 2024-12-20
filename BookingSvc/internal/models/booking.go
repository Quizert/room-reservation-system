package models

import "time"

type BookingRequest struct {
	RoomID    int       `json:"room_id"`
	HotelID   int       `json:"hotel_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`

	HotelName       string `json:"hotel_name"`
	RoomDescription string `json:"room_description"`
	RoomNumber      int    `json:"room_number"`

	CardNumber string `json:"card_number"`
	Amount     int    `json:"amount"`
}

type BookingInfo struct {
	UserID    int       `json:"user_id"`
	RoomID    int       `json:"room_id"`
	HotelID   int       `json:"hotel_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

func (req *BookingRequest) ToBookingInfo(userID int) *BookingInfo {
	return &BookingInfo{
		UserID:    userID,
		RoomID:    req.RoomID,
		HotelID:   req.HotelID,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}
}
