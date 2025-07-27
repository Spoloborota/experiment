package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Spoloborota/experiment/internal/domain/entities"
	"github.com/Spoloborota/experiment/internal/domain/repositories"
	"github.com/Spoloborota/experiment/internal/infrastructure/database/sqlc"
)

type profileRepository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

// NewProfileRepository создает новый экземпляр репозитория профилей
func NewProfileRepository(db *sql.DB) repositories.ProfileRepository {
	return &profileRepository{
		db:      db,
		queries: sqlc.New(db),
	}
}

// Create создает новый профиль
func (r *profileRepository) Create(ctx context.Context, profile *entities.Profile) (*entities.Profile, error) {
	sqlcProfile, err := r.queries.CreateProfile(ctx, sqlc.CreateProfileParams{
		UserID:    int32(profile.UserID),
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		Age:       sql.NullInt32{Int32: int32(profile.Age), Valid: true},
		Gender:    sql.NullString{String: profile.Gender, Valid: profile.Gender != ""},
		City:      sql.NullString{String: profile.City, Valid: profile.City != ""},
		Interests: profile.Interests,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create profile: %w", err)
	}

	return r.convertToEntity(sqlcProfile), nil
}

// GetByID получает профиль по ID
func (r *profileRepository) GetByID(ctx context.Context, id int) (*entities.Profile, error) {
	sqlcProfile, err := r.queries.GetProfileByID(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("profile not found")
		}
		return nil, fmt.Errorf("failed to get profile by id: %w", err)
	}

	return r.convertToEntity(sqlcProfile), nil
}

// GetByUserID получает профиль по ID пользователя
func (r *profileRepository) GetByUserID(ctx context.Context, userID int) (*entities.Profile, error) {
	sqlcProfile, err := r.queries.GetProfileByUserID(ctx, int32(userID))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("profile not found")
		}
		return nil, fmt.Errorf("failed to get profile by user id: %w", err)
	}

	return r.convertToEntity(sqlcProfile), nil
}

// Update обновляет профиль
func (r *profileRepository) Update(ctx context.Context, profile *entities.Profile) (*entities.Profile, error) {
	sqlcProfile, err := r.queries.UpdateProfile(ctx, sqlc.UpdateProfileParams{
		UserID:    int32(profile.UserID),
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		Age:       sql.NullInt32{Int32: int32(profile.Age), Valid: true},
		Gender:    sql.NullString{String: profile.Gender, Valid: profile.Gender != ""},
		City:      sql.NullString{String: profile.City, Valid: profile.City != ""},
		Interests: profile.Interests,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	return r.convertToEntity(sqlcProfile), nil
}

// Search ищет профили по фильтрам
func (r *profileRepository) Search(ctx context.Context, filters repositories.SearchFilters) ([]*entities.Profile, error) {
	var gender string
	if filters.Gender != nil {
		gender = *filters.Gender
	}

	var city string
	if filters.City != nil {
		city = *filters.City
	}

	var interests []string
	if len(filters.Interests) > 0 {
		interests = filters.Interests
	}

	sqlcProfiles, err := r.queries.SearchProfiles(ctx, sqlc.SearchProfilesParams{
		Column1: gender,
		Column2: city,
		Column3: interests,
		Limit:   int32(filters.Limit),
		Offset:  int32(filters.Offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search profiles: %w", err)
	}

	profiles := make([]*entities.Profile, len(sqlcProfiles))
	for i, sqlcProfile := range sqlcProfiles {
		profiles[i] = r.convertToEntity(sqlcProfile)
	}

	return profiles, nil
}

// Count возвращает количество профилей по фильтрам
func (r *profileRepository) Count(ctx context.Context, filters repositories.SearchFilters) (int, error) {
	var gender string
	if filters.Gender != nil {
		gender = *filters.Gender
	}

	var city string
	if filters.City != nil {
		city = *filters.City
	}

	var interests []string
	if len(filters.Interests) > 0 {
		interests = filters.Interests
	}

	count, err := r.queries.GetProfilesCount(ctx, sqlc.GetProfilesCountParams{
		Column1: gender,
		Column2: city,
		Column3: interests,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to count profiles: %w", err)
	}

	return int(count), nil
}

// convertToEntity конвертирует sqlc модель в доменную сущность
func (r *profileRepository) convertToEntity(sqlcProfile sqlc.Profile) *entities.Profile {
	var age int
	if sqlcProfile.Age.Valid {
		age = int(sqlcProfile.Age.Int32)
	}

	var gender string
	if sqlcProfile.Gender.Valid {
		gender = sqlcProfile.Gender.String
	}

	var city string
	if sqlcProfile.City.Valid {
		city = sqlcProfile.City.String
	}

	interests := []string(sqlcProfile.Interests)
	if interests == nil {
		interests = []string{}
	}

	return &entities.Profile{
		ID:        int(sqlcProfile.ID),
		UserID:    int(sqlcProfile.UserID),
		FirstName: sqlcProfile.FirstName,
		LastName:  sqlcProfile.LastName,
		Age:       age,
		Gender:    gender,
		City:      city,
		Interests: interests,
		CreatedAt: sqlcProfile.CreatedAt,
		UpdatedAt: sqlcProfile.UpdatedAt,
	}
}
