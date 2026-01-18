package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"workflow-approval/package/workflow_step/domain"
	"workflow-approval/package/workflow_step/ports"
	"workflow-approval/utils"
)

var ErrStepNotFound = errors.New("workflow step not found")

// WorkflowStepRepositoryImpl implements WorkflowStepRepository interface
type WorkflowStepRepositoryImpl struct {
	db *gorm.DB
}

// NewWorkflowStepRepository creates a new WorkflowStepRepositoryImpl instance
func NewWorkflowStepRepository(db *gorm.DB) ports.WorkflowStepRepository {
	return &WorkflowStepRepositoryImpl{db: db}
}

// Create creates a new workflow step
func (r *WorkflowStepRepositoryImpl) Create(ctx context.Context, step *domain.WorkflowStep) error {
	return r.db.WithContext(ctx).Create(step).Error
}

// GetByID retrieves a workflow step by ID
func (r *WorkflowStepRepositoryImpl) GetByID(ctx context.Context, id string) (*domain.WorkflowStep, error) {
	var step domain.WorkflowStep
	result := r.db.WithContext(ctx).First(&step, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrStepNotFound
		}
		return nil, result.Error
	}
	return &step, nil
}

// Update updates an existing workflow step
func (r *WorkflowStepRepositoryImpl) Update(ctx context.Context, step *domain.WorkflowStep) error {
	step.UpdatedAt = utils.TimeNowUTC()
	return r.db.WithContext(ctx).Save(step).Error
}

// Delete deletes a workflow step by ID
func (r *WorkflowStepRepositoryImpl) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&domain.WorkflowStep{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetByWorkflowID retrieves all steps for a workflow
func (r *WorkflowStepRepositoryImpl) GetByWorkflowID(ctx context.Context, workflowID string) ([]*domain.WorkflowStep, error) {
	var steps []*domain.WorkflowStep
	err := r.db.WithContext(ctx).
		Where("workflow_id = ?", workflowID).
		Order("level ASC").
		Find(&steps).Error
	return steps, err
}

// GetByWorkflowAndLevel retrieves a step by workflow ID and level
func (r *WorkflowStepRepositoryImpl) GetByWorkflowAndLevel(ctx context.Context, workflowID string, level int) (*domain.WorkflowStep, error) {
	var step domain.WorkflowStep
	result := r.db.WithContext(ctx).
		Where("workflow_id = ? AND level = ?", workflowID, level).
		First(&step)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrStepNotFound
		}
		return nil, result.Error
	}
	return &step, nil
}
