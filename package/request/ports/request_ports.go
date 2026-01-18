package ports

import (
	"context"

	"workflow-approval/package/request/domain"
)

// RequestRepository defines the interface for request data access
//
//go:generate mockery --with-expecter --name=RequestRepository --output=mocks --filename=RequestRepository.go
type RequestRepository interface {
	Create(ctx context.Context, request *domain.Request) error
	GetByID(ctx context.Context, id string) (*domain.Request, error)
	Update(ctx context.Context, request *domain.Request) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, limit int, status *domain.RequestStatus) ([]*domain.Request, int64, error)
	GetByIDForUpdate(ctx context.Context, id string) (*domain.Request, error) // For transaction locking
}

// RequestService defines the interface for request business logic
//
//go:generate mockery --with-expecter --name=RequestService --output=mocks --filename=RequestService.go
type RequestService interface {
	CreateRequest(ctx context.Context, workflowID, requesterID string, amount float64, title, description string) (*domain.Request, error)
	GetRequest(ctx context.Context, id string) (*domain.Request, error)
	ListRequests(ctx context.Context, page, limit int, status *domain.RequestStatus) ([]*domain.Request, int64, error)
	Approve(ctx context.Context, requestID, userID, actorID string, isAdmin bool) (*domain.Request, error)
	Reject(ctx context.Context, requestID, userID, actorID string, isAdmin bool, reason string) (*domain.Request, error)
	UpdateRequest(ctx context.Context, id string, amount float64, title, description string) (*domain.Request, error)
	DeleteRequest(ctx context.Context, id string) error

	// LockRequest acquires a mutex lock for the given request ID to prevent concurrent approval operations.
	// Returns a function that must be called to release the lock (defer it).
	LockRequest(requestID string) func()
}
