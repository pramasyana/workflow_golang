package usecase

import (
	"context"
	"errors"

	actorPorts "workflow-approval/package/actor/ports"
	actorRepo "workflow-approval/package/actor/repository"
	workflowPorts "workflow-approval/package/workflow/ports"
	workflowRepo "workflow-approval/package/workflow/repository"
	stepDomain "workflow-approval/package/workflow_step/domain"
	stepPorts "workflow-approval/package/workflow_step/ports"
	stepRepo "workflow-approval/package/workflow_step/repository"
)

var (
	ErrStepLevelRequired = errors.New("step level is required")
	ErrStepActorRequired = errors.New("step actor is required")
	ErrStepLevelExists   = errors.New("step level already exists for this workflow")
	ErrActorNotFound     = errors.New("actor not found")
	ErrWorkflowNotFound  = errors.New("workflow not found")
)

// WorkflowStepServiceImpl implements WorkflowStepService interface
type WorkflowStepServiceImpl struct {
	stepRepo     stepPorts.WorkflowStepRepository
	actorRepo    actorPorts.ActorRepository
	workflowRepo workflowPorts.WorkflowRepository
}

// NewWorkflowStepService creates a new WorkflowStepServiceImpl instance
func NewWorkflowStepService(StepRepo stepPorts.WorkflowStepRepository, ActorRepo actorPorts.ActorRepository, WorkflowRepo workflowPorts.WorkflowRepository) stepPorts.WorkflowStepService {
	return &WorkflowStepServiceImpl{
		stepRepo:     StepRepo,
		actorRepo:    ActorRepo,
		workflowRepo: WorkflowRepo,
	}
}

// CreateStep creates a new workflow step
func (s *WorkflowStepServiceImpl) CreateStep(ctx context.Context, workflowID string, level int, actorID string, conditions stepDomain.StepConditions) (*stepDomain.WorkflowStep, error) {
	if workflowID == "" {
		return nil, ErrWorkflowNotFound
	}
	if level < 1 {
		return nil, ErrStepLevelRequired
	}
	if actorID == "" {
		return nil, ErrStepActorRequired
	}

	// Validate workflow exists in database
	_, err := s.workflowRepo.GetByID(ctx, workflowID)
	if err != nil {
		if errors.Is(err, workflowRepo.ErrWorkflowNotFound) {
			return nil, ErrWorkflowNotFound
		}
		return nil, err
	}

	// Validate actor exists in database
	_, err = s.actorRepo.GetByID(ctx, actorID)
	if err != nil {
		if errors.Is(err, actorRepo.ErrActorNotFound) {
			return nil, ErrActorNotFound
		}
		return nil, err
	}

	// Check if level already exists for this workflow
	existing, err := s.stepRepo.GetByWorkflowAndLevel(ctx, workflowID, level)
	if err == nil && existing != nil {
		return nil, ErrStepLevelExists
	}
	if err != nil && !errors.Is(err, stepRepo.ErrStepNotFound) {
		return nil, err
	}

	step := stepDomain.NewWorkflowStep(workflowID, level, actorID, conditions)
	if err := s.stepRepo.Create(ctx, step); err != nil {
		return nil, err
	}

	return step, nil
}

// GetSteps retrieves all steps for a workflow
func (s *WorkflowStepServiceImpl) GetSteps(ctx context.Context, workflowID string) ([]*stepDomain.WorkflowStep, error) {
	if workflowID == "" {
		return nil, ErrWorkflowNotFound
	}

	// Validate workflow exists in database
	_, err := s.workflowRepo.GetByID(ctx, workflowID)
	if err != nil {
		if errors.Is(err, workflowRepo.ErrWorkflowNotFound) {
			return nil, ErrWorkflowNotFound
		}
		return nil, err
	}

	return s.stepRepo.GetByWorkflowID(ctx, workflowID)
}

// GetStepByID retrieves a step by ID
func (s *WorkflowStepServiceImpl) GetStepByID(ctx context.Context, id string) (*stepDomain.WorkflowStep, error) {
	return s.stepRepo.GetByID(ctx, id)
}

// UpdateStep updates a workflow step
func (s *WorkflowStepServiceImpl) UpdateStep(ctx context.Context, id string, level int, actorID string, conditions stepDomain.StepConditions) (*stepDomain.WorkflowStep, error) {
	if level < 1 {
		return nil, ErrStepLevelRequired
	}
	if actorID == "" {
		return nil, ErrStepActorRequired
	}

	// Validate actor exists in database
	_, err := s.actorRepo.GetByID(ctx, actorID)
	if err != nil {
		if errors.Is(err, actorRepo.ErrActorNotFound) {
			return nil, ErrActorNotFound
		}
		return nil, err
	}

	step, err := s.stepRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validate workflow exists in database
	_, err = s.workflowRepo.GetByID(ctx, step.WorkflowID)
	if err != nil {
		if errors.Is(err, workflowRepo.ErrWorkflowNotFound) {
			return nil, ErrWorkflowNotFound
		}
		return nil, err
	}

	step.Level = level
	step.ActorID = actorID
	step.Conditions = conditions

	if err := s.stepRepo.Update(ctx, step); err != nil {
		return nil, err
	}

	return step, nil
}

// DeleteStep deletes a step by ID
func (s *WorkflowStepServiceImpl) DeleteStep(ctx context.Context, id string) error {
	step, err := s.stepRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, stepRepo.ErrStepNotFound) {
			return err
		}
		return err
	}

	// Validate workflow exists in database
	_, err = s.workflowRepo.GetByID(ctx, step.WorkflowID)
	if err != nil {
		if errors.Is(err, workflowRepo.ErrWorkflowNotFound) {
			return ErrWorkflowNotFound
		}
		return err
	}

	return s.stepRepo.Delete(ctx, id)
}
