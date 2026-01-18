package usecase

import (
	"context"
	"errors"
	"regexp"

	"workflow-approval/package/actor/domain"
	"workflow-approval/package/actor/ports"
	"workflow-approval/package/actor/repository"
)

var (
	ErrInvalidActorName  = errors.New("invalid actor name")
	ErrActorCodeRequired = errors.New("actor code is required")
	ErrActorNameRequired = errors.New("actor name is required")
	ErrActorAlreadyExist = errors.New("actor already exists")
)

// ActorServiceImpl implements ActorService interface
type ActorServiceImpl struct {
	actorRepo ports.ActorRepository
}

// NewActorService creates a new ActorServiceImpl instance
func NewActorService(actorRepo ports.ActorRepository) ports.ActorService {
	return &ActorServiceImpl{actorRepo: actorRepo}
}

// CreateActor creates a new actor
func (s *ActorServiceImpl) CreateActor(ctx context.Context, name, code string) (*domain.Actor, error) {
	// Validation
	if name == "" {
		return nil, ErrActorNameRequired
	}
	if code == "" {
		return nil, ErrActorCodeRequired
	}

	// Validate code format (alphanumeric and underscore only)
	if !isValidCode(code) {
		return nil, errors.New("actor code must be alphanumeric with underscores")
	}

	// Check if actor already exists
	_, err := s.actorRepo.GetByCode(ctx, code)
	if err == nil {
		return nil, ErrActorAlreadyExist
	}
	if !errors.Is(err, repository.ErrActorNotFound) {
		return nil, err
	}

	// Create actor
	actor := domain.NewActor(name, code)
	if err := s.actorRepo.Create(ctx, actor); err != nil {
		return nil, err
	}

	return actor, nil
}

// GetActorByID retrieves an actor by ID
func (s *ActorServiceImpl) GetActorByID(ctx context.Context, id string) (*domain.Actor, error) {
	return s.actorRepo.GetByID(ctx, id)
}

// GetAllActors retrieves all actors
func (s *ActorServiceImpl) GetAllActors(ctx context.Context) ([]*domain.Actor, error) {
	return s.actorRepo.GetAll(ctx)
}

// UpdateActor updates an actor's name
func (s *ActorServiceImpl) UpdateActor(ctx context.Context, id string, name string) (*domain.Actor, error) {
	if name == "" {
		return nil, ErrActorNameRequired
	}

	actor, err := s.actorRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	actor.Name = name
	if err := s.actorRepo.Update(ctx, actor); err != nil {
		return nil, err
	}

	return actor, nil
}

// DeleteActor deletes an actor by ID
func (s *ActorServiceImpl) DeleteActor(ctx context.Context, id string) error {
	_, err := s.actorRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.actorRepo.Delete(ctx, id)
}

// isValidCode validates actor code format (alphanumeric and underscores)
func isValidCode(code string) bool {
	pattern := `^[a-zA-Z0-9_]+$`
	matched, _ := regexp.MatchString(pattern, code)
	return matched
}
