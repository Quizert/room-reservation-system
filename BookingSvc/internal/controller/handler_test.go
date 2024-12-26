package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/Quizert/room-reservation-system/BookingSvc/internal/mocks"
	"github.com/Quizert/room-reservation-system/BookingSvc/internal/models"
	"github.com/Quizert/room-reservation-system/BookingSvc/internal/myerror"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

// TestCreateBookingHandler_Success проверяет успешное создание бронирования
func TestCreateBookingHandler_Success(t *testing.T) {
	logger, err := zap.NewDevelopment()
	assert.NoError(t, err)
	defer logger.Sync()

	tracer := otel.Tracer("test-tracer")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingService := mocks.NewMockBookingService(ctrl)

	bookingHandler := NewBookingHandler(mockBookingService, tracer)

	bookingRequest := models.BookingRequest{
		RoomID:          1,
		HotelID:         1,
		HotelName:       "Test Hotel",
		RoomDescription: "Deluxe Room",
		StartDate:       time.Now().Add(24 * time.Hour),
		EndDate:         time.Now().Add(48 * time.Hour),
		CountOfPeople:   2,
		RoomBasePrice:   100,
		CardNumber:      "4111111111111111",
	}
	body, err := json.Marshal(bookingRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/bookings", bytes.NewBuffer(body))
	ctx := context.WithValue(req.Context(), "user_id", 1)
	ctx = context.WithValue(ctx, "username", "testuser")
	ctx = context.WithValue(ctx, "chat_id", "testchat")
	req = req.WithContext(ctx)

	mockBookingService.
		EXPECT().
		CreateBooking(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, br *models.BookingRequest, user *models.User) error {
			// Можно добавить дополнительные проверки аргументов здесь
			assert.Equal(t, 1, user.UserID)
			assert.Equal(t, "testuser", user.Username)
			assert.Equal(t, "testchat", user.ChatID)
			return nil
		}).
		Times(1)

	rr := httptest.NewRecorder()

	bookingHandler.CreateBooking(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
}

// TestCreateBookingHandler_BadRequest проверяет обработку некорректного запроса
func TestCreateBookingHandler_BadRequest(t *testing.T) {

	tracer := otel.Tracer("test-tracer")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingService := mocks.NewMockBookingService(ctrl)

	bookingHandler := NewBookingHandler(mockBookingService, tracer)

	invalidJSON := `{"room_id": "invalid_id"}` // room_id должен быть int

	req := httptest.NewRequest(http.MethodPost, "/bookings", bytes.NewBufferString(invalidJSON))
	ctx := context.WithValue(req.Context(), "user_id", 1)
	ctx = context.WithValue(ctx, "username", "testuser")
	ctx = context.WithValue(ctx, "chat_id", "testchat")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	bookingHandler.CreateBooking(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// TestCreateBookingHandler_Conflict проверяет обработку случая, когда бронирование уже существует
func TestCreateBookingHandler_Conflict(t *testing.T) {
	tracer := otel.Tracer("test-tracer")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingService := mocks.NewMockBookingService(ctrl)

	bookingHandler := NewBookingHandler(mockBookingService, tracer)

	bookingRequest := models.BookingRequest{
		RoomID:          1,
		HotelID:         1,
		HotelName:       "Test Hotel",
		RoomDescription: "Deluxe Room",
		StartDate:       time.Now().Add(24 * time.Hour),
		EndDate:         time.Now().Add(48 * time.Hour),
		CountOfPeople:   2,
		RoomBasePrice:   100,
		CardNumber:      "4111111111111111",
	}
	body, err := json.Marshal(bookingRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/bookings", bytes.NewBuffer(body))
	ctx := context.WithValue(req.Context(), "user_id", 1)
	ctx = context.WithValue(ctx, "username", "testuser")
	ctx = context.WithValue(ctx, "chat_id", "testchat")
	req = req.WithContext(ctx)

	mockBookingService.
		EXPECT().
		CreateBooking(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(myerror.ErrBookingAlreadyExists).
		Times(1)

	rr := httptest.NewRecorder()

	bookingHandler.CreateBooking(rr, req)

	assert.Equal(t, http.StatusConflict, rr.Code)
	assert.Equal(t, myerror.ErrBookingAlreadyExists.Error()+"\n", rr.Body.String())
}

// TestCreateBookingHandler_InternalServerError проверяет обработку внутренней ошибки сервера
func TestCreateBookingHandler_InternalServerError(t *testing.T) {
	tracer := otel.Tracer("test-tracer")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingService := mocks.NewMockBookingService(ctrl)

	bookingHandler := NewBookingHandler(mockBookingService, tracer)

	bookingRequest := models.BookingRequest{
		RoomID:          1,
		HotelID:         1,
		HotelName:       "Test Hotel",
		RoomDescription: "Deluxe Room",
		StartDate:       time.Now().Add(24 * time.Hour),
		EndDate:         time.Now().Add(48 * time.Hour),
		CountOfPeople:   2,
		RoomBasePrice:   100,
		CardNumber:      "4111111111111111",
	}
	body, err := json.Marshal(bookingRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/bookings", bytes.NewBuffer(body))
	ctx := context.WithValue(req.Context(), "user_id", 1)
	ctx = context.WithValue(ctx, "username", "testuser")
	ctx = context.WithValue(ctx, "chat_id", "testchat")
	req = req.WithContext(ctx)

	mockBookingService.
		EXPECT().
		CreateBooking(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(errors.New("error")).
		Times(1)

	rr := httptest.NewRecorder()

	bookingHandler.CreateBooking(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "server error"+"\n", rr.Body.String())
}

func TestGetBookingByUserID_Success(t *testing.T) {
	tracer := otel.Tracer("test-tracer")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingService := mocks.NewMockBookingService(ctrl)
	bookingHandler := NewBookingHandler(mockBookingService, tracer)

	userID := 1
	bookings := []*models.BookingInfo{
		{
			UserID:  userID,
			RoomID:  101,
			HotelID: 1,
		},
	}

	mockBookingService.
		EXPECT().
		GetBookingsByUserID(gomock.Any(), userID).
		Return(bookings, nil)

	req := httptest.NewRequest(http.MethodGet, "/bookings/users?user_id="+strconv.Itoa(userID), nil)
	req = req.WithContext(createContext(req.Context(), userID))
	rr := httptest.NewRecorder()

	bookingHandler.GetBookingByUserID(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	var result []*models.BookingInfo
	assert.NoError(t, json.NewDecoder(rr.Body).Decode(&result))
	assert.Equal(t, bookings, result)
}

func TestGetBookingByUserID_InvalidUserID(t *testing.T) {
	tracer := otel.Tracer("test-tracer")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingService := mocks.NewMockBookingService(ctrl)
	bookingHandler := NewBookingHandler(mockBookingService, tracer)

	req := httptest.NewRequest(http.MethodGet, "/bookings/users?user_id=invalid", nil)
	req = req.WithContext(createContext(req.Context(), 1))
	rr := httptest.NewRecorder()

	bookingHandler.GetBookingByUserID(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "Invalid user id\n", rr.Body.String())
}

func TestGetBookingByUserID_ForbiddenAccess(t *testing.T) {
	tracer := otel.Tracer("test-tracer")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingService := mocks.NewMockBookingService(ctrl)
	bookingHandler := NewBookingHandler(mockBookingService, tracer)

	req := httptest.NewRequest(http.MethodGet, "/bookings/users?user_id=2", nil)
	req = req.WithContext(createContext(req.Context(), 1))
	rr := httptest.NewRecorder()

	bookingHandler.GetBookingByUserID(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
	assert.Equal(t, "forbidden access\n", rr.Body.String())
}

func TestGetBookingByUserID_InternalServerError(t *testing.T) {
	tracer := otel.Tracer("test-tracer")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingService := mocks.NewMockBookingService(ctrl)
	bookingHandler := NewBookingHandler(mockBookingService, tracer)

	userID := 1

	mockBookingService.
		EXPECT().
		GetBookingsByUserID(gomock.Any(), userID).
		Return(nil, errors.New("database error"))

	req := httptest.NewRequest(http.MethodGet, "/bookings/users?user_id="+strconv.Itoa(userID), nil)
	req = req.WithContext(createContext(req.Context(), userID))
	rr := httptest.NewRecorder()

	bookingHandler.GetBookingByUserID(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "server error\n", rr.Body.String())
}

func TestGetBookingByHotelID_Success(t *testing.T) {
	tracer := otel.Tracer("test-tracer")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingService := mocks.NewMockBookingService(ctrl)
	bookingHandler := NewBookingHandler(mockBookingService, tracer)

	userID := 1
	hotelID := 1
	bookings := []*models.BookingInfo{
		{
			UserID:  userID,
			RoomID:  101,
			HotelID: hotelID,
		},
	}

	mockBookingService.
		EXPECT().
		GetBookingsByHotelID(gomock.Any(), hotelID, userID).
		Return(bookings, nil)

	req := httptest.NewRequest(http.MethodGet, "/bookings/hotels?hotel_id="+strconv.Itoa(hotelID), nil)
	req = req.WithContext(createContext(req.Context(), userID))
	rr := httptest.NewRecorder()

	bookingHandler.GetBookingByHotelID(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	var result []*models.BookingInfo
	assert.NoError(t, json.NewDecoder(rr.Body).Decode(&result))
	assert.Equal(t, bookings, result)
}

func TestGetBookingByHotelID_InvalidHotelID(t *testing.T) {
	tracer := otel.Tracer("test-tracer")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingService := mocks.NewMockBookingService(ctrl)
	bookingHandler := NewBookingHandler(mockBookingService, tracer)

	req := httptest.NewRequest(http.MethodGet, "/bookings/hotels?hotel_id=invalid", nil)
	req = req.WithContext(createContext(req.Context(), 1))
	rr := httptest.NewRecorder()

	bookingHandler.GetBookingByHotelID(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "Invalid user id\n", rr.Body.String())
}

func TestGetBookingByHotelID_ForbiddenAccess(t *testing.T) {
	tracer := otel.Tracer("test-tracer")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingService := mocks.NewMockBookingService(ctrl)
	bookingHandler := NewBookingHandler(mockBookingService, tracer)

	req := httptest.NewRequest(http.MethodGet, "/bookings/hotels?hotel_id=1", nil)
	req = req.WithContext(createContext(req.Context(), 2)) // userID does not match owner ID
	rr := httptest.NewRecorder()

	mockBookingService.
		EXPECT().
		GetBookingsByHotelID(gomock.Any(), 1, 2).
		Return(nil, myerror.ErrForbiddenAccess)

	bookingHandler.GetBookingByHotelID(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
	assert.Equal(t, "forbidden access\n", rr.Body.String())
}

func TestGetBookingByHotelID_HotelNotFound(t *testing.T) {
	tracer := otel.Tracer("test-tracer")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingService := mocks.NewMockBookingService(ctrl)
	bookingHandler := NewBookingHandler(mockBookingService, tracer)

	req := httptest.NewRequest(http.MethodGet, "/bookings/hotels?hotel_id=1", nil)
	req = req.WithContext(createContext(req.Context(), 1))
	rr := httptest.NewRecorder()

	mockBookingService.
		EXPECT().
		GetBookingsByHotelID(gomock.Any(), 1, 1).
		Return(nil, myerror.ErrHotelNotFound)

	bookingHandler.GetBookingByHotelID(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, "hotel not found\n", rr.Body.String())
}

func TestGetBookingByHotelID_InternalServerError(t *testing.T) {
	tracer := otel.Tracer("test-tracer")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingService := mocks.NewMockBookingService(ctrl)
	bookingHandler := NewBookingHandler(mockBookingService, tracer)

	userID := 1
	hotelID := 1

	mockBookingService.
		EXPECT().
		GetBookingsByHotelID(gomock.Any(), hotelID, userID).
		Return(nil, errors.New("database error"))

	req := httptest.NewRequest(http.MethodGet, "/bookings/hotels?hotel_id="+strconv.Itoa(hotelID), nil)
	req = req.WithContext(createContext(req.Context(), userID))
	rr := httptest.NewRecorder()

	bookingHandler.GetBookingByHotelID(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "server error\n", rr.Body.String())
}

//func TestGetAvailableRooms_Success(t *testing.T) {
//	tracer := otel.Tracer("test-tracer")
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockBookingService := mocks.NewMockBookingService(ctrl)
//	bookingHandler := NewBookingHandler(mockBookingService, tracer)
//
//	hotelID := 1
//	startDate := time.Now().Add(24 * time.Hour).UTC()
//	endDate := time.Now().Add(48 * time.Hour).UTC()
//	rooms := []*hotelpb.Room{
//		{
//			Id:          101,
//			HotelId: int32(hotelID),
//			RoomNumber:  "101A",
//			Description: "Deluxe Room",
//		},
//	}
//
//	mockBookingService.
//		EXPECT().
//		GetAvailableRooms(gomock.Any(), hotelID, startDate, endDate).
//		Return(rooms, nil)
//
//	req := httptest.NewRequest(http.MethodGet, "/bookings/hotels/rooms?hotel_id="+strconv.Itoa(hotelID)+"&start_date="+startDate.Format(time.RFC3339)+"&end_date="+endDate.Format(time.RFC3339), nil)
//	rr := httptest.NewRecorder()
//
//	bookingHandler.GetAvailableRooms(rr, req)
//
//	assert.Equal(t, http.StatusOK, rr.Code)
//	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
//	var result []*models.Room
//	assert.NoError(t, json.NewDecoder(rr.Body).Decode(&result))
//	assert.Equal(t, rooms, result)
//}
//
//func TestGetAvailableRooms_InvalidHotelID(t *testing.T) {
//	tracer := otel.Tracer("test-tracer")
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockBookingService := mocks.NewMockBookingService(ctrl)
//	bookingHandler := NewBookingHandler(mockBookingService, tracer)
//
//	req := httptest.NewRequest(http.MethodGet, "/bookings/hotels/rooms?hotel_id=invalid&start_date=2023-01-01T00:00:00Z&end_date=2023-01-02T00:00:00Z", nil)
//	rr := httptest.NewRecorder()
//
//	bookingHandler.GetAvailableRooms(rr, req)
//
//	assert.Equal(t, http.StatusBadRequest, rr.Code)
//	assert.Equal(t, "Invalid hotel_id\n", rr.Body.String())
//}
//
//func TestGetAvailableRooms_InvalidStartDate(t *testing.T) {
//	tracer := otel.Tracer("test-tracer")
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockBookingService := mocks.NewMockBookingService(ctrl)
//	bookingHandler := NewBookingHandler(mockBookingService, tracer)
//
//	req := httptest.NewRequest(http.MethodGet, "/bookings/hotels/rooms?hotel_id=1&start_date=invalid&end_date=2023-01-02T00:00:00Z", nil)
//	rr := httptest.NewRecorder()
//
//	bookingHandler.GetAvailableRooms(rr, req)
//
//	assert.Equal(t, http.StatusBadRequest, rr.Code)
//	assert.Contains(t, rr.Body.String(), "Invalid start_date")
//}
//
//func TestGetAvailableRooms_InvalidEndDate(t *testing.T) {
//	tracer := otel.Tracer("test-tracer")
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockBookingService := mocks.NewMockBookingService(ctrl)
//	bookingHandler := NewBookingHandler(mockBookingService, tracer)
//
//	req := httptest.NewRequest(http.MethodGet, "/bookings/hotels/rooms?hotel_id=1&start_date=2023-01-01T00:00:00Z&end_date=invalid", nil)
//	rr := httptest.NewRecorder()
//
//	bookingHandler.GetAvailableRooms(rr, req)
//
//	assert.Equal(t, http.StatusBadRequest, rr.Code)
//	assert.Contains(t, rr.Body.String(), "Invalid end_date")
//}
//
//func TestGetAvailableRooms_InternalServerError(t *testing.T) {
//	tracer := otel.Tracer("test-tracer")
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockBookingService := mocks.NewMockBookingService(ctrl)
//	bookingHandler := NewBookingHandler(mockBookingService, tracer)
//
//	hotelID := 1
//	startDate := time.Now().Add(24 * time.Hour).UTC()
//	endDate := time.Now().Add(48 * time.Hour).UTC()
//
//	mockBookingService.
//		EXPECT().
//		GetAvailableRooms(gomock.Any(), hotelID, startDate, endDate).
//		Return(nil, errors.New("database error"))
//
//	req := httptest.NewRequest(http.MethodGet, "/bookings/hotels/rooms?hotel_id="+strconv.Itoa(hotelID)+"&start_date="+startDate.Format(time.RFC3339)+"&end_date="+endDate.Format(time.RFC3339), nil)
//	rr := httptest.NewRecorder()
//
//	bookingHandler.GetAvailableRooms(rr, req)
//
//	assert.Equal(t, http.StatusInternalServerError, rr.Code)
//	assert.Equal(t, "server error\n", rr.Body.String())
//}

func createContext(ctx context.Context, userID int) context.Context {
	ctx = context.WithValue(ctx, "user_id", userID)
	ctx = context.WithValue(ctx, "username", "testuser")
	ctx = context.WithValue(ctx, "chat_id", "testchat")
	return ctx
}
