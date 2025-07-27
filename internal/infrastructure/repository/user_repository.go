package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Spoloborota/experiment/internal/domain/entities"
	"github.com/Spoloborota/experiment/internal/domain/repositories"
	"github.com/Spoloborota/experiment/internal/infrastructure/database/sqlc"
)

type userRepository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

// NewUserRepository создает новый экземпляр репозитория пользователей
func NewUserRepository(db *sql.DB) repositories.UserRepository {
	return &userRepository{
		db:      db,
		queries: sqlc.New(db),
	}
}

// Create создает нового пользователя
func (r *userRepository) Create(ctx context.Context, user *entities.User) (*entities.User, error) {
	sqlcUser, err := r.queries.CreateUser(ctx, sqlc.CreateUserParams{
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return r.convertToEntity(sqlcUser), nil
}

// GetByID получает пользователя по ID
func (r *userRepository) GetByID(ctx context.Context, id int) (*entities.User, error) {
	sqlcUser, err := r.queries.GetUserByID(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return r.convertToEntity(sqlcUser), nil
}

// GetByEmail получает пользователя по email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	sqlcUser, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return r.convertToEntity(sqlcUser), nil
}

// convertToEntity конвертирует sqlc модель в доменную сущность
func (r *userRepository) convertToEntity(sqlcUser sqlc.User) *entities.User {
	return &entities.User{
		ID:           int(sqlcUser.ID),
		Email:        sqlcUser.Email,
		PasswordHash: sqlcUser.PasswordHash,
		CreatedAt:    sqlcUser.CreatedAt,
		UpdatedAt:    sqlcUser.UpdatedAt,
	}
}
