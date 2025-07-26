package entities

import (
	"errors"
	"strings"
	"time"
)

type Profile struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Age       int       `json:"age"`
	Gender    string    `json:"gender"`
	City      string    `json:"city"`
	Interests []string  `json:"interests"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

// NewProfile создает новый профиль с валидацией
func NewProfile(userID int, firstName, lastName string, age int, gender string, city string, interests []string) (*Profile, error) {
	if err := validateProfileData(firstName, lastName, age, gender); err != nil {
		return nil, err
	}

	// Очищаем и нормализуем данные
	cleanInterests := cleanInterests(interests)

	return &Profile{
		UserID:    userID,
		FirstName: strings.TrimSpace(firstName),
		LastName:  strings.TrimSpace(lastName),
		Age:       age,
		Gender:    strings.ToLower(gender),
		City:      strings.TrimSpace(city),
		Interests: cleanInterests,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// Update обновляет профиль с валидацией
func (p *Profile) Update(firstName, lastName string, age int, gender string, city string, interests []string) error {
	if err := validateProfileData(firstName, lastName, age, gender); err != nil {
		return err
	}

	p.FirstName = strings.TrimSpace(firstName)
	p.LastName = strings.TrimSpace(lastName)
	p.Age = age
	p.Gender = strings.ToLower(gender)
	p.City = strings.TrimSpace(city)
	p.Interests = cleanInterests(interests)
	p.UpdatedAt = time.Now()

	return nil
}

// GetFullName возвращает полное имя
func (p *Profile) GetFullName() string {
	return p.FirstName + " " + p.LastName
}

// validateProfileData проверяет корректность данных профиля
func validateProfileData(firstName, lastName string, age int, gender string) error {
	if strings.TrimSpace(firstName) == "" {
		return errors.New("first name cannot be empty")
	}

	if strings.TrimSpace(lastName) == "" {
		return errors.New("last name cannot be empty")
	}

	if age < 16 || age > 120 {
		return errors.New("age must be between 16 and 120")
	}

	validGenders := map[string]bool{
		string(GenderMale):   true,
		string(GenderFemale): true,
		string(GenderOther):  true,
	}

	if !validGenders[strings.ToLower(gender)] {
		return errors.New("invalid gender value")
	}

	return nil
}

// cleanInterests очищает и нормализует список интересов
func cleanInterests(interests []string) []string {
	var cleaned []string
	seen := make(map[string]bool)

	for _, interest := range interests {
		trimmed := strings.TrimSpace(interest)
		if trimmed != "" && !seen[trimmed] {
			cleaned = append(cleaned, trimmed)
			seen[trimmed] = true
		}
	}

	return cleaned
}
