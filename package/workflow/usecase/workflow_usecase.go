package usecase

import (
	"context"
	"errors"

	"workflow-approval/package/workflow/domain"
	"workflow-approval/package/workflow/ports"
	"workflow-approval/package/workflow/repository"
)

var (
	ErrWorkflowNameRequired = errors.New("workflow name is required")
)

// WorkflowServiceImpl implements WorkflowService interface
type WorkflowServiceImpl struct {
	workflowRepo ports.WorkflowRepository
}

// NewWorkflowService creates a new WorkflowServiceImpl instance
func NewWorkflowService(workflowRepo ports.WorkflowRepository) ports.WorkflowService {
	return &WorkflowServiceImpl{workflowRepo: workflowRepo}
}

// CreateWorkflow creates a new workflow
func (s *WorkflowServiceImpl) CreateWorkflow(ctx context.Context, name string) (*domain.Workflow, error) {
	if name == "" {
		return nil, ErrWorkflowNameRequired
	}

	workflow := domain.NewWorkflow(name)
	if err := s.workflowRepo.Create(ctx, workflow); err != nil {
		return nil, err
	}

	return workflow, nil
}

// GetWorkflow retrieves a workflow by ID
func (s *WorkflowServiceImpl) GetWorkflow(ctx context.Context, id string) (*domain.Workflow, error) {
	return s.workflowRepo.GetByID(ctx, id)
}

// UpdateWorkflow updates a workflow's name
func (s *WorkflowServiceImpl) UpdateWorkflow(ctx context.Context, id string, name string) (*domain.Workflow, error) {
	if name == "" {
		return nil, ErrWorkflowNameRequired
	}

	workflow, err := s.workflowRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	workflow.Name = name
	if err := s.workflowRepo.Update(ctx, workflow); err != nil {
		return nil, err
	}

	return workflow, nil
}

// ListWorkflows retrieves a paginated list of workflows
func (s *WorkflowServiceImpl) ListWorkflows(ctx context.Context, page, limit int) ([]*domain.Workflow, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return s.workflowRepo.List(ctx, page, limit)
}

// DeleteWorkflow deletes a workflow by ID
func (s *WorkflowServiceImpl) DeleteWorkflow(ctx context.Context, id string) error {
	_, err := s.workflowRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrWorkflowNotFound) {
			return err
		}
		return err
	}
	return s.workflowRepo.Delete(ctx, id)
}
