package ports

import (
	"context"

	actorDomain "workflow-approval/package/actor/domain"
	userDomain "workflow-approval/package/user/domain"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(ctx context.Context, user *userDomain.User) error
	GetByID(ctx context.Context, id string) (*userDomain.User, error)
	GetByEmail(ctx context.Context, email string) (*userDomain.User, error)
	Update(ctx context.Context, user *userDomain.User) error
	Delete(ctx context.Context, id string) error
}

// ActorRepository defines the interface for actor data access
type ActorRepository interface {
	GetByID(ctx context.Context, id string) (*actorDomain.Actor, error)
}

// UserService defines the interface for user business logic
type UserService interface {
	Register(ctx context.Context, email, password, name string, isAdmin bool, actorID *string) (*userDomain.User, error)
	Login(ctx context.Context, email, password string) (*userDomain.User, string, error)
	GetUserByID(ctx context.Context, id string) (*userDomain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*userDomain.User, error)
	UpdateUser(ctx context.Context, id string, name string) (*userDomain.User, error)
}
