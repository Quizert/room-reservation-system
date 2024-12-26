package models

type BookingMessage struct {
	BookingID       int    `json:"booking_id"`
	HotelID         int    `json:"hotel_id"`
	HotelName       string `json:"hotel_name"`
	RoomDescription string `json:"room_description"`
	RoomNumber      int    `json:"room_number"`
	Username        string `json:"user_name"`
	ChatID          string `json:"chat_id"`
	StartDate       string `json:"start_date"`
	EndDate         string `json:"end_date"`
}

func (req *BookingRequest) ToBookingMessage(bookingID int, username, chatID, startDate, endDate string) *BookingMessage {
	return &BookingMessage{
		BookingID:       bookingID,
		HotelID:         req.HotelID,
		HotelName:       req.HotelName,
		RoomDescription: req.RoomDescription,
		RoomNumber:      req.RoomNumber,
		StartDate:       startDate,
		EndDate:         endDate,
		Username:        username,
		ChatID:          chatID,
	}
}

func (message *BookingMessage) ToHotelierMessage(hotelierName string, hotelierChatID string) *BookingMessage {
	return &BookingMessage{
		BookingID:       message.BookingID,
		HotelName:       message.HotelName,
		RoomDescription: message.RoomDescription,
		RoomNumber:      message.RoomNumber,
		StartDate:       message.StartDate,
		EndDate:         message.EndDate,
		Username:        hotelierName,
		ChatID:          hotelierChatID,
	}
}
