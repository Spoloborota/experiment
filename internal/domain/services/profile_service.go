package services

import (
	"context"
	"errors"

	"github.com/Spoloborota/experiment/internal/domain/entities"
	"github.com/Spoloborota/experiment/internal/domain/repositories"
)

type ProfileService struct {
	profileRepo repositories.ProfileRepository
}

func NewProfileService(profileRepo repositories.ProfileRepository) *ProfileService {
	return &ProfileService{
		profileRepo: profileRepo,
	}
}

// CreateProfile создает новый профиль
func (s *ProfileService) CreateProfile(ctx context.Context, userID int, firstName, lastName string, age int, gender, city string, interests []string) (*entities.Profile, error) {
	// Проверяем, есть ли уже профиль у пользователя
	existingProfile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err == nil && existingProfile != nil {
		return nil, errors.New("profile already exists for this user")
	}

	// Создаем новый профиль
	profile, err := entities.NewProfile(userID, firstName, lastName, age, gender, city, interests)
	if err != nil {
		return nil, err
	}

	// Сохраняем в базе
	return s.profileRepo.Create(ctx, profile)
}

// GetProfile получает профиль по ID
func (s *ProfileService) GetProfile(ctx context.Context, id int) (*entities.Profile, error) {
	return s.profileRepo.GetByID(ctx, id)
}

// GetProfileByUserID получает профиль по ID пользователя
func (s *ProfileService) GetProfileByUserID(ctx context.Context, userID int) (*entities.Profile, error) {
	return s.profileRepo.GetByUserID(ctx, userID)
}

// UpdateProfile обновляет профиль
func (s *ProfileService) UpdateProfile(ctx context.Context, userID int, firstName, lastName string, age int, gender, city string, interests []string) (*entities.Profile, error) {
	// Получаем существующий профиль
	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("profile not found")
	}

	// Обновляем данные
	err = profile.Update(firstName, lastName, age, gender, city, interests)
	if err != nil {
		return nil, err
	}

	// Сохраняем изменения
	return s.profileRepo.Update(ctx, profile)
}

// SearchProfiles ищет профили по фильтрам
func (s *ProfileService) SearchProfiles(ctx context.Context, filters repositories.SearchFilters) ([]*entities.Profile, int, error) {
	// Устанавливаем значения по умолчанию
	if filters.Limit <= 0 {
		filters.Limit = 10
	}
	if filters.Limit > 100 {
		filters.Limit = 100 // Ограничиваем максимальное количество
	}
	if filters.Offset < 0 {
		filters.Offset = 0
	}

	// Получаем профили
	profiles, err := s.profileRepo.Search(ctx, filters)
	if err != nil {
		return nil, 0, err
	}

	// Получаем общее количество
	total, err := s.profileRepo.Count(ctx, filters)
	if err != nil {
		return nil, 0, err
	}

	return profiles, total, nil
}
