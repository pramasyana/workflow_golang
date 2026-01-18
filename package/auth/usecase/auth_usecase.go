package usecase

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"workflow-approval/package/auth/ports"
	"workflow-approval/package/user/domain"
	"workflow-approval/utils/jwthelper"
)

var (
	ErrInvalidLogin = errors.New("invalid email or password")
	ErrInvalidToken = errors.New("invalid token")
	ErrUserNotFound = errors.New("user not found")
)

// AuthServiceImpl implements AuthService interface
type AuthServiceImpl struct {
	authRepo  ports.AuthRepository
	jwtHelper *jwthelper.JWTHelper
}

// NewAuthService creates a new AuthServiceImpl instance
func NewAuthService(authRepo ports.AuthRepository, jwtHelper *jwthelper.JWTHelper) ports.AuthService {
	return &AuthServiceImpl{
		authRepo:  authRepo,
		jwtHelper: jwtHelper,
	}
}

// Login authenticates a user and returns a JWT token
func (s *AuthServiceImpl) Login(ctx context.Context, email, password string) (*domain.User, string, error) {
	// Validation
	if email == "" || password == "" {
		return nil, "", ErrInvalidLogin
	}

	// Get user by email
	user, err := s.authRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, "", ErrInvalidLogin
		}
		return nil, "", err
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", ErrInvalidLogin
	}

	return user, "", nil
}

// Refresh generates a new JWT token for an existing user
func (s *AuthServiceImpl) Refresh(ctx context.Context, userID string) (*domain.User, string, error) {
	// Get user by ID
	user, err := s.authRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, "", ErrUserNotFound
		}
		return nil, "", err
	}

	// Generate new JWT token
	token, err := s.jwtHelper.GenerateJWT(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// Logout validates a token (placeholder for future token blacklist implementation)
func (s *AuthServiceImpl) Logout(ctx context.Context, token string) error {
	// Validate token
	_, err := s.jwtHelper.ValidateToken(token)
	if err != nil {
		return ErrInvalidToken
	}

	// In a production system, you would add the token to a blacklist
	// For now, we just validate and acknowledge the logout
	return nil
}
