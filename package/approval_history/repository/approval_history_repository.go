package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"workflow-approval/package/approval_history/domain"
	"workflow-approval/package/approval_history/ports"
)

// ApprovalHistoryRepositoryImpl implements ApprovalHistoryRepository interface
type ApprovalHistoryRepositoryImpl struct {
	db *gorm.DB
}

// NewApprovalHistoryRepository creates a new ApprovalHistoryRepositoryImpl instance
func NewApprovalHistoryRepository(db *gorm.DB) ports.ApprovalHistoryRepository {
	return &ApprovalHistoryRepositoryImpl{db: db}
}

// Create creates a new approval history entry
func (r *ApprovalHistoryRepositoryImpl) Create(ctx context.Context, history *domain.ApprovalHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

// GetByID retrieves an approval history entry by ID
func (r *ApprovalHistoryRepositoryImpl) GetByID(ctx context.Context, id string) (*domain.ApprovalHistory, error) {
	var history domain.ApprovalHistory
	err := r.db.WithContext(ctx).First(&history, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrApprovalHistoryNotFound
		}
		return nil, err
	}
	return &history, nil
}

// GetByRequestID retrieves all approval history entries for a request
func (r *ApprovalHistoryRepositoryImpl) GetByRequestID(ctx context.Context, requestID string) ([]*domain.ApprovalHistory, error) {
	var histories []*domain.ApprovalHistory
	err := r.db.WithContext(ctx).
		Where("request_id = ?", requestID).
		Order("created_at DESC").
		Find(&histories).Error
	return histories, err
}

// GetByActorID retrieves all approval history entries by an actor
func (r *ApprovalHistoryRepositoryImpl) GetByActorID(ctx context.Context, actorID string) ([]*domain.ApprovalHistory, error) {
	var histories []*domain.ApprovalHistory
	err := r.db.WithContext(ctx).
		Where("actor_id = ?", actorID).
		Order("created_at DESC").
		Find(&histories).Error
	return histories, err
}

// GetByRequestIDOrdered retrieves all approval history for a request ordered by created_at ascending
func (r *ApprovalHistoryRepositoryImpl) GetByRequestIDOrdered(ctx context.Context, requestID string) ([]*domain.ApprovalHistory, error) {
	var histories []*domain.ApprovalHistory
	err := r.db.WithContext(ctx).
		Where("request_id = ?", requestID).
		Order("created_at ASC").
		Find(&histories).Error
	return histories, err
}

// ErrApprovalHistoryNotFound is returned when the approval history is not found
var ErrApprovalHistoryNotFound = gorm.ErrRecordNotFound
