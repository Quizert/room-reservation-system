package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Quizert/room-reservation-system/BookingSvc/internal/models"
	"github.com/Quizert/room-reservation-system/BookingSvc/internal/myerror"
	"github.com/Quizert/room-reservation-system/HotelSvc/api/grpc/hotelpb"
	"github.com/Quizert/room-reservation-system/Libs/metrics"
	"log"
	"net/http"
	"strconv"
	"time"
)

type BookingService interface {
	CreateBooking(ctx context.Context, bookingRequest *models.BookingRequest, user *models.User) error
	GetBookingsByUserID(ctx context.Context, userID int) ([]*models.BookingInfo, error)
	GetBookingsByHotelID(ctx context.Context, hotelID, userID int) ([]*models.BookingInfo, error)
	GetAvailableRooms(ctx context.Context, hotelID int, startDate, endDate time.Time) ([]*hotelpb.Room, error)
	UpdateBookingStatus(ctx context.Context, status string, bookingMessage *models.BookingMessage) error
}

type BookingHandler struct {
	bookingService BookingService
}

func NewBookingHandler(b BookingService) *BookingHandler {
	return &BookingHandler{b}
}

func (b *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	status := http.StatusCreated
	defer func() {
		duration := time.Since(start).Seconds()
		metrics.RecordHttpMetrics(r.Method, "/bookings", http.StatusText(status), duration)
	}()
	ctx := r.Context()

	userID := ctx.Value("user_id").(int)
	username := ctx.Value("username").(string)
	chatID := ctx.Value("chat_id").(string)

	user := models.NewUser(userID, username, chatID)

	var bookingRequest models.BookingRequest
	if err := json.NewDecoder(r.Body).Decode(&bookingRequest); err != nil {
		status = http.StatusBadRequest
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := b.bookingService.CreateBooking(ctx, &bookingRequest, user); err != nil {
		if errors.Is(err, myerror.ErrBookingAlreadyExists) {
			status = http.StatusConflict
			http.Error(w, myerror.ErrBookingAlreadyExists.Error(), http.StatusConflict)
			return
		}
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (b *BookingHandler) GetBookingByUserID(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	status := http.StatusOK
	defer func() {
		duration := time.Since(start).Seconds()
		metrics.RecordHttpMetrics(r.Method, "/bookings/users", http.StatusText(status), duration)
	}()
	ctx := r.Context()

	userID := ctx.Value("user_id").(int)

	userIDParams, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		status = http.StatusBadRequest
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}
	if userIDParams != userID {
		status = http.StatusForbidden
		http.Error(w, "forbidden access", http.StatusForbidden)
		return
	}

	bookings, err := b.bookingService.GetBookingsByUserID(ctx, userID)
	if err != nil {
		status = http.StatusInternalServerError
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}

func (b *BookingHandler) GetBookingByHotelID(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	status := http.StatusOK
	defer func() {
		duration := time.Since(start).Seconds()
		metrics.RecordHttpMetrics(r.Method, "/bookings/hotels", http.StatusText(status), duration)
	}()
	ctx := r.Context()

	userID := ctx.Value("user_id").(int) // Должен быть Владелец отеля
	hotelID, err := strconv.Atoi(r.URL.Query().Get("hotel_id"))
	if err != nil {
		status = http.StatusBadRequest
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	bookings, err := b.bookingService.GetBookingsByHotelID(ctx, hotelID, userID)
	if err != nil {
		if errors.Is(err, myerror.ErrForbiddenAccess) {
			status = http.StatusForbidden
			http.Error(w, "forbidden access", http.StatusForbidden)
			return
		} else if errors.Is(err, myerror.ErrHotelNotFound) {
			status = http.StatusNotFound
			http.Error(w, "hotel not found", http.StatusNotFound)
			return
		}
		status = http.StatusInternalServerError
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}

func (b *BookingHandler) GetAvailableRooms(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	status := http.StatusOK
	defer func() {
		duration := time.Since(start).Seconds()
		metrics.RecordHttpMetrics(r.Method, "/bookings/hotels/rooms", http.StatusText(status), duration)
	}()

	ctx := r.Context()
	hotelId, err := strconv.Atoi(r.URL.Query().Get("hotel_id"))
	if err != nil {
		status = http.StatusBadRequest
		http.Error(w, "Invalid hotel_id", http.StatusBadRequest)
		return
	}
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	// Парсим start_date
	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		status = http.StatusBadRequest
		http.Error(w, fmt.Sprintf("Invalid start_date: %v", err), http.StatusBadRequest)
		return
	}
	// Парсим end_date
	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		status = http.StatusBadRequest
		http.Error(w, fmt.Sprintf("Invalid end_date: %v", err), http.StatusBadRequest)
		return
	}
	availableRooms, err := b.bookingService.GetAvailableRooms(ctx, hotelId, startDate.UTC(), endDate.UTC())
	if err != nil {
		status = http.StatusInternalServerError
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(availableRooms)

}

func (b *BookingHandler) HandlePaymentWebHook(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	status := http.StatusOK
	defer func() {
		duration := time.Since(start).Seconds()
		metrics.RecordHttpMetrics(r.Method, "/bookings/payment/response", http.StatusText(status), duration)
	}()
	ctx := r.Context()
	var paymentResponse models.PaymentResponse
	if err := json.NewDecoder(r.Body).Decode(&paymentResponse); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err := b.bookingService.UpdateBookingStatus(ctx, paymentResponse.Status, paymentResponse.MetaData)

	switch paymentResponse.Status {
	case "success":
		if err != nil {
			status = http.StatusInternalServerError
			log.Println("handler UpdateBookingStatusSuccess: ", err.Error())
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success booking!"))
	case "failed":

	}
}
