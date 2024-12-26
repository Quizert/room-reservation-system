package models

type User struct {
	ID         int
	Username   string `json:"username"`
	ChatID     string `json:"chat_id"`
	Password   string `json:"password"`
	IsHotelier bool   `json:"is_hotelier"`
}
