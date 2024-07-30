package handlers

import (
	"net/http"
	"net/http/httptest"
	"social_network/internal/models"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByName(ctx context.Context, firstName, lastName string) (*models.User, error) {
	args := m.Called(ctx, firstName, lastName)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.User), args.Error(1)
}

func TestRegister(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/register", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	userRepo := new(MockUserRepository)
	userRepo.On("CreateUser", mock.Anything, mock.Anything).Return(nil)

	h := NewAuthHandler(userRepo)

	if assert.NoError(t, h.Register(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

func TestLogin(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/login", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	userRepo := new(MockUserRepository)
	user := &models.User{
		FirstName:    "John",
		LastName:     "Doe",
		PasswordHash: "$2a$10$yH.ZpO6z/8KUgH0HgPSgZuH/JFYPx5R.LLoP.0MCTCfK9zYwJu/gS", // hashed password
	}
	userRepo.On("GetUserByName", mock.Anything, user.FirstName, user.LastName).Return(user, nil)

	h := NewAuthHandler(userRepo)

	if assert.NoError(t, h.Login(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
