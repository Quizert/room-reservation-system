package models

// Hotel представляет данные об отеле.
type Hotel struct {
	Id          int    `json:"id"`
	OwnerId     int    `json:"OwnerId"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
