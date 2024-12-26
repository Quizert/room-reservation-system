package postgresql

import (
	"database/sql"
	"errors"
	"github.com/Quizert/room-reservation-system/HotelSvc/internal/models"
)

type PostgresHotelRepository struct {
	db *sql.DB
}

func NewPostgresHotelRepository(db *sql.DB) *PostgresHotelRepository {
	return &PostgresHotelRepository{db: db}
}

func (repo *PostgresHotelRepository) GetAllHotels() ([]models.Hotel, error) {
	rows, err := repo.db.Query("SELECT Id, OwnerId, Name, Description FROM hotels")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hotels []models.Hotel
	for rows.Next() {
		var hotel models.Hotel
		if err := rows.Scan(&hotel.Id, &hotel.OwnerId, &hotel.Name, &hotel.Description); err != nil {
			return nil, err
		}
		hotels = append(hotels, hotel)
	}

	return hotels, rows.Err()
}

func (repo *PostgresHotelRepository) AddHotel(hotel models.Hotel) error {
	_, err := repo.db.Exec(
		"INSERT INTO hotels (OwnerId, Name, Description) VALUES ($1, $2, $3)",
		hotel.OwnerId, hotel.Name, hotel.Description,
	)
	return err
}

func (repo *PostgresHotelRepository) UpdateHotel(hotel models.Hotel) error {
	result, err := repo.db.Exec(
		"UPDATE hotels SET name = $1, description = $2 WHERE id = $3",
		hotel.Name, hotel.Description, hotel.Id,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("hotel not found")
	}
	return nil
}

func (repo *PostgresHotelRepository) GetHotelByID(id int) (*models.Hotel, error) {
	var hotel models.Hotel
	err := repo.db.QueryRow("SELECT Id, OwnerId, Name, Description FROM hotels WHERE id = $1", id).
		Scan(&hotel.Id, &hotel.OwnerId, &hotel.Name, &hotel.Description)
	if err == sql.ErrNoRows {
		return nil, errors.New("hotel not found")
	}
	return &hotel, err
}
