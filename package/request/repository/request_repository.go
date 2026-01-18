package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"workflow-approval/package/request/domain"
	"workflow-approval/package/request/ports"
	"workflow-approval/utils"
)

var ErrRequestNotFound = errors.New("request not found")

// RequestRepositoryImpl implements RequestRepository interface with optimistic locking
type RequestRepositoryImpl struct {
	db *gorm.DB
}

// NewRequestRepository creates a new RequestRepositoryImpl instance
func NewRequestRepository(db *gorm.DB) ports.RequestRepository {
	return &RequestRepositoryImpl{db: db}
}

// Create creates a new request
func (r *RequestRepositoryImpl) Create(ctx context.Context, request *domain.Request) error {
	return r.db.WithContext(ctx).Create(request).Error
}

// GetByID retrieves a request by ID
func (r *RequestRepositoryImpl) GetByID(ctx context.Context, id string) (*domain.Request, error) {
	var request domain.Request
	result := r.db.WithContext(ctx).First(&request, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrRequestNotFound
		}
		return nil, result.Error
	}
	return &request, nil
}

// GetByIDForUpdate retrieves a request with row-level locking for update
// This is used within a transaction to prevent race conditions
func (r *RequestRepositoryImpl) GetByIDForUpdate(ctx context.Context, id string) (*domain.Request, error) {
	var request domain.Request
	result := r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&request, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrRequestNotFound
		}
		return nil, result.Error
	}
	return &request, nil
}

// Update updates a request with optimistic locking
// Uses version field to prevent concurrent updates
func (r *RequestRepositoryImpl) Update(ctx context.Context, request *domain.Request) error {
	request.UpdatedAt = utils.TimeNowUTC()

	result := r.db.WithContext(ctx).
		Where("id = ? AND version = ?", request.ID, request.Version).
		Save(request)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("optimistic lock error: request was modified by another transaction")
	}

	return nil
}

// Delete deletes a request by ID
func (r *RequestRepositoryImpl) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&domain.Request{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// List retrieves a paginated list of requests
func (r *RequestRepositoryImpl) List(ctx context.Context, page, limit int, status *domain.RequestStatus) ([]*domain.Request, int64, error) {
	var requests []*domain.Request
	var total int64

	offset := (page - 1) * limit

	query := r.db.WithContext(ctx).Model(&domain.Request{})

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&requests).Error; err != nil {
		return nil, 0, err
	}

	return requests, total, nil
}
