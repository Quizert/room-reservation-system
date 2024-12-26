package service

//
//import (
//	"context"
//	"github.com/Quizert/room-reservation-system/BookingSvc/internal/models"
//	"github.com/Quizert/room-reservation-system/BookingSvc/internal/service/mocks"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/mock"
//	"testing"
//)
//
//func TestCreateBooking_Success(t *testing.T) {
//	// Arrange
//	storageMock := &mocks.MockStorage{}
//	paymentMock := &mocks.MockPaymentSystemClient{}
//
//	// Устанавливаем ожидание для вызова CreateBooking с любыми аргументами
//	storageMock.EXPECT().CreateBooking(mock.Anything, mock.Anything).Return(123, nil).Times(1)
//
//	// Устанавливаем ожидание для вызова CreatePaymentRequest с любыми аргументами
//	paymentMock.EXPECT().CreatePaymentRequest(mock.Anything, mock.Anything).Return(nil).Times(1)
//
//	bookingService := &BookingService{
//		storage:             storageMock,
//		paymentSystemClient: paymentMock,
//	}
//
//	bookingRequest := &models.BookingRequest{
//		CardNumber: "1234",
//		Amount:     1000,
//	}
//
//	// Act
//	err := bookingService.CreateBooking(context.Background(), bookingRequest)
//
//	// Assert
//	assert.NoError(t, err)
//}
