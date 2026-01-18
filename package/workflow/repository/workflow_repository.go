package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"workflow-approval/package/workflow/domain"
	"workflow-approval/package/workflow/ports"
	"workflow-approval/utils"
)

var ErrWorkflowNotFound = errors.New("workflow not found")

// WorkflowRepositoryImpl implements WorkflowRepository interface
type WorkflowRepositoryImpl struct {
	db *gorm.DB
}

// NewWorkflowRepository creates a new WorkflowRepositoryImpl instance
func NewWorkflowRepository(db *gorm.DB) ports.WorkflowRepository {
	return &WorkflowRepositoryImpl{db: db}
}

// Create creates a new workflow
func (r *WorkflowRepositoryImpl) Create(ctx context.Context, workflow *domain.Workflow) error {
	return r.db.WithContext(ctx).Create(workflow).Error
}

// GetByID retrieves a workflow by ID
func (r *WorkflowRepositoryImpl) GetByID(ctx context.Context, id string) (*domain.Workflow, error) {
	var workflow domain.Workflow
	result := r.db.WithContext(ctx).First(&workflow, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrWorkflowNotFound
		}
		return nil, result.Error
	}
	return &workflow, nil
}

// Update updates an existing workflow
func (r *WorkflowRepositoryImpl) Update(ctx context.Context, workflow *domain.Workflow) error {
	workflow.UpdatedAt = utils.TimeNowUTC()
	return r.db.WithContext(ctx).Save(workflow).Error
}

// Delete deletes a workflow by ID
func (r *WorkflowRepositoryImpl) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&domain.Workflow{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// List retrieves a paginated list of workflows
func (r *WorkflowRepositoryImpl) List(ctx context.Context, page, limit int) ([]*domain.Workflow, int64, error) {
	var workflows []*domain.Workflow
	var total int64

	offset := (page - 1) * limit

	// Get total count
	if err := r.db.WithContext(ctx).Model(&domain.Workflow{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&workflows).Error; err != nil {
		return nil, 0, err
	}

	return workflows, total, nil
}
