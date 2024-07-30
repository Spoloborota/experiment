package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"social_network/internal/models"
	"testing"
)

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userRepo := NewUserRepository(db)

	user := &models.User{
		FirstName:    "John",
		LastName:     "Doe",
		Age:          30,
		Gender:       "Male",
		Interests:    "Reading",
		City:         "New York",
		PasswordHash: "hashed_password",
	}

	mock.ExpectQuery(`INSERT INTO users`).
		WithArgs(user.FirstName, user.LastName, user.Age, user.Gender, user.Interests, user.City, user.PasswordHash).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	err = userRepo.CreateUser(context.Background(), user)
	assert.NoError(t, err)
	assert.Equal(t, 1, user.ID)
}
