package ports

import (
	"context"

	"workflow-approval/package/workflow/domain"
)

// WorkflowRepository defines the interface for workflow data access
//
//go:generate mockery --with-expecter --name=WorkflowRepository --output=mocks --filename=WorkflowRepository.go
type WorkflowRepository interface {
	Create(ctx context.Context, workflow *domain.Workflow) error
	GetByID(ctx context.Context, id string) (*domain.Workflow, error)
	Update(ctx context.Context, workflow *domain.Workflow) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, limit int) ([]*domain.Workflow, int64, error)
}

// WorkflowService defines the interface for workflow business logic
//
//go:generate mockery --with-expecter --name=WorkflowService --output=mocks --filename=WorkflowService.go
type WorkflowService interface {
	CreateWorkflow(ctx context.Context, name string) (*domain.Workflow, error)
	GetWorkflow(ctx context.Context, id string) (*domain.Workflow, error)
	UpdateWorkflow(ctx context.Context, id string, name string) (*domain.Workflow, error)
	ListWorkflows(ctx context.Context, page, limit int) ([]*domain.Workflow, int64, error)
	DeleteWorkflow(ctx context.Context, id string) error
}
