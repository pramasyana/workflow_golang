package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"workflow-approval/package/actor/domain"
	"workflow-approval/package/actor/ports"
)

var (
	ErrActorNotFound     = errors.New("actor not found")
	ErrActorAlreadyExist = errors.New("actor already exists")
)

// ActorRepositoryImpl implements ActorRepository interface
type ActorRepositoryImpl struct {
	db *gorm.DB
}

// NewActorRepository creates a new ActorRepositoryImpl instance
func NewActorRepository(db *gorm.DB) ports.ActorRepository {
	return &ActorRepositoryImpl{db: db}
}

// Create creates a new actor
func (r *ActorRepositoryImpl) Create(ctx context.Context, actor *domain.Actor) error {
	result := r.db.WithContext(ctx).Create(actor)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return ErrActorAlreadyExist
		}
		return result.Error
	}
	return nil
}

// GetByID retrieves an actor by ID
func (r *ActorRepositoryImpl) GetByID(ctx context.Context, id string) (*domain.Actor, error) {
	var actor domain.Actor
	result := r.db.WithContext(ctx).First(&actor, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrActorNotFound
		}
		return nil, result.Error
	}
	return &actor, nil
}

// GetByCode retrieves an actor by code
func (r *ActorRepositoryImpl) GetByCode(ctx context.Context, code string) (*domain.Actor, error) {
	var actor domain.Actor
	result := r.db.WithContext(ctx).First(&actor, "code = ?", code)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrActorNotFound
		}
		return nil, result.Error
	}
	return &actor, nil
}

// GetAll retrieves all actors
func (r *ActorRepositoryImpl) GetAll(ctx context.Context) ([]*domain.Actor, error) {
	var actors []*domain.Actor
	result := r.db.WithContext(ctx).Find(&actors)
	if result.Error != nil {
		return nil, result.Error
	}
	return actors, nil
}

// Update updates an existing actor
func (r *ActorRepositoryImpl) Update(ctx context.Context, actor *domain.Actor) error {
	result := r.db.WithContext(ctx).Save(actor)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Delete deletes an actor by ID
func (r *ActorRepositoryImpl) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&domain.Actor{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
