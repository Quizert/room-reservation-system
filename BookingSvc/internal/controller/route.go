package controller

import (
	"github.com/Quizert/room-reservation-system/Libs/middleware"
	"net/http"
)

func SetupRoutes(bookingHandler *BookingHandler) *http.ServeMux {
	mux := http.NewServeMux()

	middlewareHandler := middleware.NewMiddleware("LUIGI")
	mux.HandleFunc("/bookings", middlewareHandler.Auth(bookingHandler.CreateBooking, false))             // POST - Создается новое бронирование
	mux.HandleFunc("/bookings/users", middlewareHandler.Auth(bookingHandler.GetBookingByUserID, false))  // GET - получаем все бронирования пользователя
	mux.HandleFunc("/bookings/hotels", middlewareHandler.Auth(bookingHandler.GetBookingByHotelID, true)) // Get - получаем все бронирования отельера

	mux.HandleFunc("/bookings/hotels/rooms", bookingHandler.GetAvailableRooms) //Тут добавить сортировку по времени
	mux.HandleFunc("/bookings/payment/response", bookingHandler.HandlePaymentWebHook)
	return mux
}
