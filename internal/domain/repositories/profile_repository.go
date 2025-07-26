package repositories

import (
	"context"

	"github.com/Spoloborota/experiment/internal/domain/entities"
)

// SearchFilters определяет фильтры для поиска профилей
type SearchFilters struct {
	Gender    *string
	City      *string
	Interests []string
	Limit     int
	Offset    int
}

// ProfileRepository определяет интерфейс для работы с профилями
type ProfileRepository interface {
	// Create создает новый профиль
	Create(ctx context.Context, profile *entities.Profile) (*entities.Profile, error)

	// GetByID получает профиль по ID
	GetByID(ctx context.Context, id int) (*entities.Profile, error)

	// GetByUserID получает профиль по ID пользователя
	GetByUserID(ctx context.Context, userID int) (*entities.Profile, error)

	// Update обновляет профиль
	Update(ctx context.Context, profile *entities.Profile) (*entities.Profile, error)

	// Search ищет профили по фильтрам
	Search(ctx context.Context, filters SearchFilters) ([]*entities.Profile, error)

	// Count возвращает количество профилей по фильтрам
	Count(ctx context.Context, filters SearchFilters) (int, error)
}
