package models

// Room представляет данные о комнате в отеле.
type Room struct {
	ID         int `json:"id"`
	HotelID    int `json:"hotel_id"`
	RoomTypeID int `json:"room_type_id"`
	Number     int `json:"number"`
}

// RoomType описывает тип комнаты и её базовую стоимость.
type RoomType struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	BasePrice   int    `json:"base_price"`
}
