package handler

import (
	"github.com/Quizert/room-reservation-system/HotelSvc/internal/service"
	"github.com/Quizert/room-reservation-system/Libs/middleware"
	"net/http"
)

func RegisterHotelRoutes(mux *http.ServeMux, hotelService *service.HotelService) {
	handler := &HotelHandler{hotelService: hotelService}

	middlewareHandler := middleware.NewMiddleware("LUIGI")
	mux.HandleFunc("/hotels", handler.GetHotels)                                       // GET - список отелей
	mux.HandleFunc("/add_hotel", middlewareHandler.Auth(handler.AddHotel, true))       // POST - добавление отеля
	mux.HandleFunc("/update_hotel", middlewareHandler.Auth(handler.UpdateHotel, true)) // PUT - обновление отеля
}
