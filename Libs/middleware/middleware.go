package middleware

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
)

type Middleware struct {
	secret string
}

func NewMiddleware(secret string) *Middleware {
	return &Middleware{secret: secret}
}

func (m *Middleware) Auth(next http.HandlerFunc, clientHotelierAccess bool) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization missing", http.StatusUnauthorized)
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}
		tokenEncoded := parts[1]
		token, err := jwt.Parse(tokenEncoded, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("alg format is wrong %v", token.Header["alg"])
			}
			return []byte(m.secret), nil
		})
		if err != nil {
			http.Error(w, "invalid auth", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID, ok := claims["user_id"].(float64)
			if !ok {
				http.Error(w, "invalid user id", http.StatusUnauthorized)
				return
			}
			isHotelier, _ := claims["is_hotelier"].(bool)
			if isHotelier != clientHotelierAccess {
				http.Error(w, "forbidden access", http.StatusForbidden)
				return
			}
			username, _ := claims["username"].(string)
			chatID, _ := claims["chat_id"].(string)
			log.Println(token)
			ctx := context.WithValue(r.Context(), "user_id", int(userID))
			ctx = context.WithValue(ctx, "is_hotelier", isHotelier) //TODO: Мейби не нужно
			ctx = context.WithValue(ctx, "username", username)
			ctx = context.WithValue(ctx, "chat_id", chatID)
			next(w, r.WithContext(ctx))
		}
	})
}
