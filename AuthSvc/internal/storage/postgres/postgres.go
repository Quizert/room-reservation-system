package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/Quizert/room-reservation-system/AuthSvc/internal/models"
	"github.com/Quizert/room-reservation-system/AuthSvc/internal/storage"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) RegisterUser(ctx context.Context, user *models.User) (int, error) {
	log.Println("--------------------", user.ChatID, "---------------------")
	query := `
		SELECT EXISTS (
			SELECT ID 
			FROM users
			WHERE ChatID = $1
		);
	`
	var exists bool
	err := r.db.QueryRow(ctx, query, user.ChatID).Scan(&exists)
	if err != nil {
		return 0, fmt.Errorf("error checking if user exists: %w", err)
	}
	if exists {
		return 0, fmt.Errorf("in register user: %w", storage.ErrUserExists)
	}

	query = `
		INSERT INTO users (Username, ChatID, Password, IsHotelier)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`

	var id int
	err = r.db.QueryRow(ctx, query, user.Username, user.ChatID, user.Password, user.IsHotelier).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error inserting user: %w", err)
	}
	return id, nil
}

func (r *Repository) LoginUser(ctx context.Context, chatID string) (*models.User, error) {
	query := `
		SELECT ID, Username, ChatID, Password, IsHotelier FROM users 
		WHERE ChatID = $1
	`

	var user models.User

	err := r.db.QueryRow(ctx, query, chatID).Scan(
		&user.ID,
		&user.Username,
		&user.ChatID,
		&user.Password,
		&user.IsHotelier,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("in login User: %w", storage.ErrUserNotFound)
		}
		return nil, fmt.Errorf("error login user: %w", err)
	}
	return &user, nil
}

func (r *Repository) IsHotelier(ctx context.Context, userID int) (bool, error) {
	query := `
		SELECT IsHotelier FROM users
		where UserID = $1
	`
	var isHotelier bool
	err := r.db.QueryRow(ctx, query, userID).Scan(&isHotelier)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, fmt.Errorf("in isHotelier User: %w", storage.ErrUserNotFound)
		}
		return false, fmt.Errorf("error isHotelier user: %w", err)
	}
	return isHotelier, nil
}
