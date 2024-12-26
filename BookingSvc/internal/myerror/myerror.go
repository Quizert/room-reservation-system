package myerror

import "errors"

var (
	ErrForbiddenAccess      = errors.New("forbidden access")
	ErrHotelNotFound        = errors.New("hotel Not Found")
	ErrBookingAlreadyExists = errors.New("booking already exists")
)
