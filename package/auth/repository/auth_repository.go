package repository

import (
	"context"
	"errors"

	"workflow-approval/package/auth/ports"
	"workflow-approval/package/user/domain"
	userPorts "workflow-approval/package/user/ports"
)

var ErrUserNotFound = errors.New("user not found")

// AuthRepositoryImpl implements AuthRepository interface
// It wraps the user repository for authentication purposes
type AuthRepositoryImpl struct {
	userRepo userPorts.UserRepository
}

// NewAuthRepository creates a new AuthRepositoryImpl instance
func NewAuthRepository(userRepo userPorts.UserRepository) ports.AuthRepository {
	return &AuthRepositoryImpl{userRepo: userRepo}
}

// GetByEmail retrieves a user by email
func (r *AuthRepositoryImpl) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return r.userRepo.GetByEmail(ctx, email)
}

// GetByID retrieves a user by ID
func (r *AuthRepositoryImpl) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return r.userRepo.GetByID(ctx, id)
}
