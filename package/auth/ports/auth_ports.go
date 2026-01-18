package ports

import (
	"context"

	"workflow-approval/package/user/domain"
)

// AuthRepository defines the interface for authentication data access
type AuthRepository interface {
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
}

// AuthService defines the interface for authentication business logic
type AuthService interface {
	Login(ctx context.Context, email, password string) (*domain.User, string, error)
	Refresh(ctx context.Context, userID string) (*domain.User, string, error)
	Logout(ctx context.Context, token string) error
}
