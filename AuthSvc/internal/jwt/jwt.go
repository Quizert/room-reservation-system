package jwt

import (
	"github.com/Quizert/room-reservation-system/AuthSvc/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func NewToken(user *models.User, secret string, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["username"] = user.Username
	claims["chat_id"] = user.ChatID
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["is_hotelier"] = user.IsHotelier

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
