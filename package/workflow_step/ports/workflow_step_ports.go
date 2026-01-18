package ports

import (
	"context"

	stepDomain "workflow-approval/package/workflow_step/domain"
)

// WorkflowStepRepository defines the interface for workflow step data access
//
//go:generate mockery --with-expecter --name=WorkflowStepRepository --output=mocks --filename=WorkflowStepRepository.go
type WorkflowStepRepository interface {
	Create(ctx context.Context, step *stepDomain.WorkflowStep) error
	GetByID(ctx context.Context, id string) (*stepDomain.WorkflowStep, error)
	Update(ctx context.Context, step *stepDomain.WorkflowStep) error
	Delete(ctx context.Context, id string) error
	GetByWorkflowID(ctx context.Context, workflowID string) ([]*stepDomain.WorkflowStep, error)
	GetByWorkflowAndLevel(ctx context.Context, workflowID string, level int) (*stepDomain.WorkflowStep, error)
}

// WorkflowStepService defines the interface for workflow step business logic
//
//go:generate mockery --with-expecter --name=WorkflowStepService --output=mocks --filename=WorkflowStepService.go
type WorkflowStepService interface {
	CreateStep(ctx context.Context, workflowID string, level int, actorID string, conditions stepDomain.StepConditions) (*stepDomain.WorkflowStep, error)
	GetSteps(ctx context.Context, workflowID string) ([]*stepDomain.WorkflowStep, error)
	GetStepByID(ctx context.Context, id string) (*stepDomain.WorkflowStep, error)
	UpdateStep(ctx context.Context, id string, level int, actorID string, conditions stepDomain.StepConditions) (*stepDomain.WorkflowStep, error)
	DeleteStep(ctx context.Context, id string) error
}
