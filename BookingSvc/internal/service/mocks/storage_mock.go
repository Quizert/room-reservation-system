// Code generated by MockGen. DO NOT EDIT.
// Source: storage.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"
	time "time"

	models "github.com/Quizert/room-reservation-system/BookingSvc/internal/models"
	gomock "github.com/golang/mock/gomock"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// CreateBooking mocks base method.
func (m *MockStorage) CreateBooking(ctx context.Context, booking *models.Booking) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBooking", ctx, booking)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateBooking indicates an expected call of CreateBooking.
func (mr *MockStorageMockRecorder) CreateBooking(ctx, booking interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBooking", reflect.TypeOf((*MockStorage)(nil).CreateBooking), ctx, booking)
}

// DeleteBooking mocks base method.
func (m *MockStorage) DeleteBooking(ctx context.Context, bookingID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBooking", ctx, bookingID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBooking indicates an expected call of DeleteBooking.
func (mr *MockStorageMockRecorder) DeleteBooking(ctx, bookingID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBooking", reflect.TypeOf((*MockStorage)(nil).DeleteBooking), ctx, bookingID)
}

// GetBookingsByHotelID mocks base method.
func (m *MockStorage) GetBookingsByHotelID(ctx context.Context, bookingID int) (*models.Booking, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBookingsByHotelID", ctx, bookingID)
	ret0, _ := ret[0].(*models.Booking)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBookingsByHotelID indicates an expected call of GetBookingsByHotelID.
func (mr *MockStorageMockRecorder) GetBookingsByHotelID(ctx, bookingID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBookingsByHotelID", reflect.TypeOf((*MockStorage)(nil).GetBookingsByHotelID), ctx, bookingID)
}

// GetBookingsByUserID mocks base method.
func (m *MockStorage) GetBookingsByUserID(ctx context.Context, userID int) ([]*models.Booking, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBookingsByUserID", ctx, userID)
	ret0, _ := ret[0].([]*models.Booking)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBookingsByUserID indicates an expected call of GetBookingsByUserID.
func (mr *MockStorageMockRecorder) GetBookingsByUserID(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBookingsByUserID", reflect.TypeOf((*MockStorage)(nil).GetBookingsByUserID), ctx, userID)
}

// GetUnavailableRoomsByHotelId mocks base method.
func (m *MockStorage) GetUnavailableRoomsByHotelId(ctx context.Context, HotelID int, startDate, endDate time.Time) (map[int]struct{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUnavailableRoomsByHotelId", ctx, HotelID, startDate, endDate)
	ret0, _ := ret[0].(map[int]struct{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUnavailableRoomsByHotelId indicates an expected call of GetUnavailableRoomsByHotelId.
func (mr *MockStorageMockRecorder) GetUnavailableRoomsByHotelId(ctx, HotelID, startDate, endDate interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUnavailableRoomsByHotelId", reflect.TypeOf((*MockStorage)(nil).GetUnavailableRoomsByHotelId), ctx, HotelID, startDate, endDate)
}

// UpdateBookingStatus mocks base method.
func (m *MockStorage) UpdateBookingStatus(ctx context.Context, status string, bookingID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBookingStatus", ctx, status, bookingID)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateBookingStatus indicates an expected call of UpdateBookingStatus.
func (mr *MockStorageMockRecorder) UpdateBookingStatus(ctx, status, bookingID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBookingStatus", reflect.TypeOf((*MockStorage)(nil).UpdateBookingStatus), ctx, status, bookingID)
}