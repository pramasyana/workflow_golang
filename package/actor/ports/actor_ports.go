package ports

import (
	"context"

	"workflow-approval/package/actor/domain"
)

// ActorRepository defines the interface for actor data access
type ActorRepository interface {
	Create(ctx context.Context, actor *domain.Actor) error
	GetByID(ctx context.Context, id string) (*domain.Actor, error)
	GetByCode(ctx context.Context, code string) (*domain.Actor, error)
	GetAll(ctx context.Context) ([]*domain.Actor, error)
	Update(ctx context.Context, actor *domain.Actor) error
	Delete(ctx context.Context, id string) error
}

// ActorService defines the interface for actor business logic
type ActorService interface {
	CreateActor(ctx context.Context, name, code string) (*domain.Actor, error)
	GetActorByID(ctx context.Context, id string) (*domain.Actor, error)
	GetAllActors(ctx context.Context) ([]*domain.Actor, error)
	UpdateActor(ctx context.Context, id string, name string) (*domain.Actor, error)
	DeleteActor(ctx context.Context, id string) error
}
