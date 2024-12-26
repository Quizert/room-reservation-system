package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Quizert/room-reservation-system/HotelSvc/internal/myerror"
)

type PostgresOwnerRepository struct {
	db *sql.DB
}

func NewPostgresOwnerRepository(db *sql.DB) *PostgresOwnerRepository {
	return &PostgresOwnerRepository{db: db}
}

func (repo *PostgresOwnerRepository) GetOwnerIdByHotelId(ctx context.Context, hotelID int) (int, error) {
	var ownerId int
	err := repo.db.QueryRow("SELECT OwnerID from hotels WHERE id = $1 ", hotelID).Scan(&ownerId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("error getting owner id from hotels: %w", myerror.ErrHotelNotFound)
		}
		return 0, fmt.Errorf("error getting owner id from hotels: %w", err)
	}
	return ownerId, err
}
