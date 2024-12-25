package controller

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Quizert/room-reservation-system/AuthSvc/internal/models"
	"github.com/Quizert/room-reservation-system/AuthSvc/internal/myerror"
	"github.com/Quizert/room-reservation-system/AuthSvc/pkj/authpb"
	"github.com/Quizert/room-reservation-system/Libs/metrics"
	"net/http"
	"time"
)

type AuthService interface {
	RegisterUser(ctx context.Context, user *models.User) (int, error)
	LoginUser(ctx context.Context, user *models.User) (string, error)
	GetHotelierInformation(ctx context.Context, request *authpb.GetHotelierRequest) (*authpb.GetHotelierResponse, error)
}

type AuthHandler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (a *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	status := http.StatusCreated
	defer func() {
		duration := time.Since(start).Seconds()
		metrics.RecordHttpMetrics(r.Method, "/auth/register", http.StatusText(status), duration)
	}()

	ctx := r.Context()
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		status = http.StatusBadRequest
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	_, err := a.authService.RegisterUser(ctx, &user)
	if err != nil {
		if errors.Is(err, myerror.ErrUserExists) {
			status = http.StatusConflict
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}
		status = http.StatusInternalServerError
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (a *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	status := http.StatusOK
	defer func() {
		duration := time.Since(start).Seconds()
		metrics.RecordHttpMetrics(r.Method, "/auth/login", http.StatusText(status), duration)
	}()
	ctx := r.Context()
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		status = http.StatusBadRequest
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	token, err := a.authService.LoginUser(ctx, &user)
	if err != nil {
		if errors.Is(err, myerror.ErrInvalidCredentials) {
			status = http.StatusNotFound
			http.Error(w, myerror.ErrInvalidCredentials.Error(), http.StatusNotFound)
			return
		}
		status = http.StatusInternalServerError
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	jwtResponse := map[string]string{
		"access_token": token,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(jwtResponse)
	if err != nil {
		status = http.StatusInternalServerError
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
