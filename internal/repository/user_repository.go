package repository

import (
	"context"
	"social_network/internal/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByName(ctx context.Context, firstName, lastName string) (*models.User, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (first_name, last_name, age, gender, interests, city, password_hash)
              VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	return r.db.QueryRow(ctx, query, user.FirstName, user.LastName, user.Age, user.Gender, user.Interests, user.City, user.PasswordHash).Scan(&user.ID)
}

func (r *userRepository) GetUserByName(ctx context.Context, firstName, lastName string) (*models.User, error) {
	query := `SELECT id, password_hash FROM users WHERE first_name = $1 AND last_name = $2`
	row := r.db.QueryRow(ctx, query, firstName, lastName)

	user := &models.User{}
	err := row.Scan(&user.ID, &user.PasswordHash)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	query := `SELECT id, first_name, last_name, age, gender, interests, city FROM users WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	user := &models.User{}
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Age, &user.Gender, &user.Interests, &user.City)
	if err != nil {
		return nil, err
	}
	return user, nil
}
