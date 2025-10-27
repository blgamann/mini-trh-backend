package services

import (
	"errors"

	"mini-trh-backend/internal/consts"
	"mini-trh-backend/internal/logger"
	"mini-trh-backend/internal/utils"
	"mini-trh-backend/pkg/domain/entities"
	"mini-trh-backend/pkg/infrastructure/postgres/repositories"

	"go.uber.org/zap"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
)

// AuthService handles authentication operations
type AuthService struct {
	userRepo   *repositories.UserRepository
	jwtService *JWTService
}

// NewAuthService creates a new AuthService instance
func NewAuthService(userRepo *repositories.UserRepository, jwtService *JWTService) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

// Register creates a new user account
func (s *AuthService) Register(email, password, role string) (*entities.User, error) {
	// Check if user already exists
	exists, err := s.userRepo.ExistsByEmail(email)
	if err != nil {
		logger.Log.Error("Failed to check user existence", zap.Error(err))
		return nil, err
	}
	if exists {
		return nil, ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		logger.Log.Error("Failed to hash password", zap.Error(err))
		return nil, err
	}

	// Create user entity
	user := &entities.User{
		Email:    email,
		Password: hashedPassword,
		Role:     role,
	}

	// Save to database
	if err := s.userRepo.Create(user); err != nil {
		logger.Log.Error("Failed to create user", zap.Error(err))
		return nil, err
	}

	logger.Log.Info("User registered successfully", zap.String("email", email))
	return user, nil
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(email, password string) (string, *entities.User, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		logger.Log.Error("Failed to find user", zap.Error(err))
		return "", nil, err
	}
	if user == nil {
		return "", nil, ErrInvalidCredentials
	}

	// Verify password
	if !utils.CheckPasswordHash(password, user.Password) {
		return "", nil, ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := s.jwtService.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		logger.Log.Error("Failed to generate token", zap.Error(err))
		return "", nil, err
	}

	logger.Log.Info("User logged in successfully", zap.String("email", email))
	return token, user, nil
}

// CreateDefaultAdmin creates a default admin user if no users exist
func (s *AuthService) CreateDefaultAdmin(email, password string) error {
	// Check if any users exist
	count, err := s.userRepo.Count()
	if err != nil {
		return err
	}

	// If users exist, don't create default admin
	if count > 0 {
		logger.Log.Info("Users already exist, skipping default admin creation")
		return nil
	}

	// Create default admin
	_, err = s.Register(email, password, consts.RoleAdmin)
	if err != nil {
		logger.Log.Error("Failed to create default admin", zap.Error(err))
		return err
	}

	logger.Log.Info("Default admin created successfully", zap.String("email", email))
	return nil
}
