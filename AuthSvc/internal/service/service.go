package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Quizert/room-reservation-system/AuthSvc/internal/jwt"
	"github.com/Quizert/room-reservation-system/AuthSvc/internal/models"
	"github.com/Quizert/room-reservation-system/AuthSvc/internal/myerror"
	"github.com/Quizert/room-reservation-system/AuthSvc/pkj/authpb"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthServiceImpl struct {
	storage  Storage
	tokenTTl time.Duration
	secret   string
	log      *zap.Logger
}

func NewAuthServiceImpl(storage Storage, tokenTTl time.Duration, secret string, log *zap.Logger) *AuthServiceImpl {
	return &AuthServiceImpl{
		storage:  storage,
		tokenTTl: tokenTTl,
		secret:   secret,
		log:      log,
	}
}

func (a *AuthServiceImpl) RegisterUser(ctx context.Context, user *models.User) (int, error) {
	a.log.With(
		zap.String("Layer", "service: RegisterUser"),
		zap.String("username", user.Username),
		zap.Bool("is_hotelier", user.IsHotelier),
		zap.String("chat_id", user.ChatID)).Info("Received request to register user")

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		a.log.Error("failed to hash password", zap.Error(err))
		return 0, fmt.Errorf("%s: %w", "auth.RegisterUser", err)
	}

	user.Password = string(passwordHash)

	id, err := a.storage.RegisterUser(ctx, user)
	if err != nil {
		if errors.Is(err, myerror.ErrUserExists) {
			a.log.Warn("user already exists", zap.Error(err))
			return 0, fmt.Errorf("%s: %w", "auth.RegisterUser", myerror.ErrUserExists)
		}
		a.log.Error("failed to register user", zap.Error(err))

		return 0, fmt.Errorf("%s: %w", "auth.RegisterUser", err)
	}
	return id, nil

}

func (a *AuthServiceImpl) LoginUser(ctx context.Context, user *models.User) (string, error) {
	a.log.With(
		zap.String("Layer", "Auth.RegisterUser"),
		zap.String("username", user.Username),
		zap.Bool("is_hotelier", user.IsHotelier),
		zap.String("chat_id", user.ChatID)).Info("Received request to login user")

	UserExist, err := a.storage.LoginUser(ctx, user.ChatID)
	if err != nil {
		if errors.Is(err, myerror.ErrUserNotFound) {
			a.log.Warn("user not found", zap.Error(err))
			return "", fmt.Errorf("%s: %w", "auth.LoginUser", myerror.ErrInvalidCredentials)
		}

		a.log.Error("failed to get user", zap.Error(err))
		return "", fmt.Errorf("%s: %w", "auth.LoginUser", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(UserExist.Password), []byte(user.Password)); err != nil {
		a.log.Warn("invalid credentials", zap.Error(err))

		return "", fmt.Errorf("%s: %w", "auth.LoginUser", myerror.ErrInvalidCredentials)
	}

	token, err := jwt.NewToken(UserExist, a.secret, a.tokenTTl)
	if err != nil {
		a.log.Error("failed to generate token", zap.Error(err))
		return "", fmt.Errorf("%s: %w", "auth.LoginUser", err)
	}

	return token, nil
}

func (a *AuthServiceImpl) GetHotelierInformation(ctx context.Context, request *authpb.GetHotelierRequest) (*authpb.GetHotelierResponse, error) {
	response, err := a.storage.GetHotelierInformation(ctx, request)
	if err != nil {
		if errors.Is(err, myerror.ErrUserNotFound) {
			a.log.Warn("user not found", zap.Error(err))
			return nil, fmt.Errorf("%s: %w", "auth.LoginUser", myerror.ErrUserNotFound)
		}
		a.log.Error("failed to get hotelier information", zap.Error(err))
	}
	return response, nil
}

//func (a *AuthServiceImpl) IsHotelier(ctx context.Context, userID int) (bool, error) {
//	a.log.With(
//		zap.String("Layer", "Auth.IsHotelier"),
//		zap.Int("user_id", userID))
//
//	IsHotelier, err := a.storage.IsHotelier(ctx, userID)
//	if err != nil {
//		if errors.Is(err, myerror.ErrUserNotFound) {
//			a.log.Warn("user not found", zap.Error(err))
//			return false, fmt.Errorf("%s: %w", "auth.LoginUser", err)
//		}
//		a.log.Error("failed to check if user is hotelier", zap.Error(err))
//		return false, fmt.Errorf("%s: %w", "auth.LoginUser", err)
//
//	}
//	return IsHotelier, nil
//}
