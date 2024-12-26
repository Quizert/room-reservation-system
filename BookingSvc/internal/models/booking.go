package models

import "time"

type BookingRequest struct {
	RoomID          int    `json:"room_id"`
	HotelID         int    `json:"hotel_id"`
	HotelName       string `json:"hotel_name"`
	RoomDescription string `json:"room_description"`
	RoomNumber      int    `json:"room_number"`
	RoomBasePrice   int    `json:"room_base_price"`

	CardNumber    string    `json:"card_number"`
	CountOfPeople int       `json:"count_of_people"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`

	Amount int `json:"-"`
}

type BookingInfo struct {
	UserID    int       `json:"user_id"`
	RoomID    int       `json:"room_id"`
	HotelID   int       `json:"hotel_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type User struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	ChatID   string `json:"chat_id"`
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

func NewUser(userID int, username string, chatID string) *User {
	return &User{
		UserID:   userID,
		Username: username,
		ChatID:   chatID,
	}
}
