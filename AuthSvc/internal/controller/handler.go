package controller

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Quizert/room-reservation-system/AuthSvc/internal/models"
	"github.com/Quizert/room-reservation-system/AuthSvc/internal/myerror"
	"github.com/Quizert/room-reservation-system/AuthSvc/pkj/authpb"
	"github.com/Quizert/room-reservation-system/Libs/metrics"
	"go.opentelemetry.io/otel/trace"
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
	tracer      trace.Tracer
}

func NewAuthHandler(authService AuthService, trace trace.Tracer) *AuthHandler {
	return &AuthHandler{authService: authService, tracer: trace}
}

func (a *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx, span := a.tracer.Start(r.Context(), "Handler.RegisterUser")
	defer span.End()

	start := time.Now()
	status := http.StatusCreated
	defer func() {
		duration := time.Since(start).Seconds()
		metrics.RecordHttpMetrics(r.Method, "/auth/register", http.StatusText(status), duration)
	}()

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		span.RecordError(err)

		status = http.StatusBadRequest
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	_, err := a.authService.RegisterUser(ctx, &user)
	if err != nil {
		span.RecordError(err)

		if errors.Is(err, myerror.ErrUserExists) {
			status = http.StatusConflict
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}
		status = http.StatusInternalServerError
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	span.AddEvent("Register user success")
	w.WriteHeader(http.StatusCreated)
}

func (a *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx, span := a.tracer.Start(r.Context(), "Handler.LoginUser")
	defer span.End()

	start := time.Now()
	status := http.StatusOK
	defer func() {
		duration := time.Since(start).Seconds()
		metrics.RecordHttpMetrics(r.Method, "/auth/login", http.StatusText(status), duration)
	}()
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		span.RecordError(err)
		status = http.StatusBadRequest
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	token, err := a.authService.LoginUser(ctx, &user)
	if err != nil {
		span.RecordError(err)
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
		span.RecordError(err)

		status = http.StatusInternalServerError
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	span.AddEvent("Login user success")
}
