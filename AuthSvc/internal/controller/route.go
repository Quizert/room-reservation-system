package controller

import "net/http"

func SetupRoutes(authHandler *AuthHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/auth/register", authHandler.RegisterUser)
	mux.HandleFunc("/auth/login", authHandler.LoginUser)
	return mux
}
