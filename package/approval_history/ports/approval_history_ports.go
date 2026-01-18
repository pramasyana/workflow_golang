package ports

import (
	"context"

	"workflow-approval/package/approval_history/domain"
	"workflow-approval/package/approval_history/domain/dto"
)

// ApprovalHistoryRepository defines the interface for approval history data access
type ApprovalHistoryRepository interface {
	// Create creates a new approval history entry
	Create(ctx context.Context, history *domain.ApprovalHistory) error

	// GetByID retrieves an approval history entry by ID
	GetByID(ctx context.Context, id string) (*domain.ApprovalHistory, error)

	// GetByRequestID retrieves all approval history entries for a request
	GetByRequestID(ctx context.Context, requestID string) ([]*domain.ApprovalHistory, error)

	// GetByActorID retrieves all approval history entries by an actor
	GetByActorID(ctx context.Context, actorID string) ([]*domain.ApprovalHistory, error)

	// GetByRequestIDOrdered retrieves all approval history for a request ordered by created_at ascending
	GetByRequestIDOrdered(ctx context.Context, requestID string) ([]*domain.ApprovalHistory, error)
}

// ApprovalHistoryService defines the interface for approval history business logic
type ApprovalHistoryService interface {
	// GetHistoryByRequestID retrieves all approval history for a request with details
	GetHistoryByRequestID(ctx context.Context, requestID string) ([]*dto.ApprovalHistoryResponse, error)

	// CreateApprovalHistory creates a new approval history entry
	CreateApprovalHistory(ctx context.Context, requestID, workflowID string, stepLevel int, actorID, userID string, action domain.ApprovalAction, comment string) error
}
