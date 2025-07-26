package repositories

import (
	"context"

	"github.com/Spoloborota/experiment/internal/domain/entities"
)

// UserRepository определяет интерфейс для работы с пользователями
type UserRepository interface {
	// Create создает нового пользователя
	Create(ctx context.Context, user *entities.User) (*entities.User, error)

	// GetByID получает пользователя по ID
	GetByID(ctx context.Context, id int) (*entities.User, error)

	// GetByEmail получает пользователя по email
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
}
