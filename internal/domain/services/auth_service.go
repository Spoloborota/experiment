package services

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/Spoloborota/experiment/internal/domain/entities"
	"github.com/Spoloborota/experiment/internal/domain/repositories"
)

type AuthService struct {
	userRepo       repositories.UserRepository
	jwtSecret      string
	jwtExpiryHours int
}

type JWTClaims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func NewAuthService(userRepo repositories.UserRepository, jwtSecret string, jwtExpiryHours int) *AuthService {
	return &AuthService{
		userRepo:       userRepo,
		jwtSecret:      jwtSecret,
		jwtExpiryHours: jwtExpiryHours,
	}
}

// Register регистрирует нового пользователя
func (s *AuthService) Register(ctx context.Context, email, password string) (*entities.User, error) {
	// Проверяем, существует ли пользователь
	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Хешируем пароль
	passwordHash, err := s.hashPassword(password)
	if err != nil {
		return nil, err
	}

	// Создаем пользователя
	user, err := entities.NewUser(email, passwordHash)
	if err != nil {
		return nil, err
	}

	// Сохраняем в базе
	return s.userRepo.Create(ctx, user)
}

// Login авторизует пользователя и возвращает JWT токен
func (s *AuthService) Login(ctx context.Context, email, password string) (string, *entities.User, error) {
	// Получаем пользователя
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	// Проверяем пароль
	if !s.checkPassword(password, user.PasswordHash) {
		return "", nil, errors.New("invalid credentials")
	}

	// Генерируем JWT токен
	token, err := s.generateJWT(user)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

// ValidateToken валидирует JWT токен и возвращает информацию о пользователе
func (s *AuthService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// GetUserByID получает пользователя по ID
func (s *AuthService) GetUserByID(ctx context.Context, userID int) (*entities.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

// hashPassword хеширует пароль с помощью bcrypt
func (s *AuthService) hashPassword(password string) (string, error) {
	if len(password) < 6 {
		return "", errors.New("password must be at least 6 characters long")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// checkPassword проверяет пароль против хеша
func (s *AuthService) checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// generateJWT генерирует JWT токен для пользователя
func (s *AuthService) generateJWT(user *entities.User) (string, error) {
	claims := &JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.jwtExpiryHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
