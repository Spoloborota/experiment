package entities

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Не включаем в JSON
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// NewUser создает нового пользователя с валидацией
func NewUser(email, passwordHash string) (*User, error) {
	if err := validateEmail(email); err != nil {
		return nil, err
	}

	if passwordHash == "" {
		return nil, errors.New("password hash cannot be empty")
	}

	return &User{
		Email:        strings.ToLower(email),
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

// validateEmail проверяет корректность email адреса
func validateEmail(email string) error {
	if email == "" {
		return errors.New("email cannot be empty")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}

	return nil
}

// TouchUpdatedAt обновляет время последнего изменения
func (u *User) TouchUpdatedAt() {
	u.UpdatedAt = time.Now()
}
