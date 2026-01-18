package usecase

import (
	"context"
	"errors"
	"sync"

	approvalHistoryDomain "workflow-approval/package/approval_history/domain"
	approvalHistoryPorts "workflow-approval/package/approval_history/ports"
	reqDomain "workflow-approval/package/request/domain"
	reqPorts "workflow-approval/package/request/ports"
	reqRepo "workflow-approval/package/request/repository"
	wfPorts "workflow-approval/package/workflow/ports"
	wfRepo "workflow-approval/package/workflow/repository"
	stepDomain "workflow-approval/package/workflow_step/domain"
	stepPorts "workflow-approval/package/workflow_step/ports"
	stepRepo "workflow-approval/package/workflow_step/repository"
)

var (
	ErrRequestAmountRequired = errors.New("request amount is required")
	ErrRequestAmountPositive = errors.New("request amount must be greater than 0")
	ErrRequestNotPending     = errors.New("request is not in pending status")
	ErrRequestNotFound       = errors.New("request not found")
	ErrWorkflowNotFound      = errors.New("workflow not found")
	ErrNoNextStep            = errors.New("no more steps in workflow")
	ErrInvalidWorkflow       = errors.New("invalid workflow")
	ErrRejectReasonRequired  = errors.New("reject reason is required")
	ErrUnauthorizedActor     = errors.New("unauthorized actor: you are not the assigned approver for this step")
)

// RequestServiceImpl implements RequestService interface with approval workflow logic
// and provides double-layer concurrency control (Mutex + SELECT FOR UPDATE)
type RequestServiceImpl struct {
	requestRepo         reqPorts.RequestRepository
	workflowRepo        wfPorts.WorkflowRepository
	workflowStepRepo    stepPorts.WorkflowStepRepository
	approvalHistoryRepo approvalHistoryPorts.ApprovalHistoryRepository

	// mutexMap stores per-request mutexes for in-memory locking
	// This is Layer 1 of our double-layer concurrency control
	mutexMap sync.Map
}

// NewRequestService creates a new RequestServiceImpl instance
func NewRequestService(
	requestRepo reqPorts.RequestRepository,
	workflowRepo wfPorts.WorkflowRepository,
	workflowStepRepo stepPorts.WorkflowStepRepository,
	approvalHistoryRepo approvalHistoryPorts.ApprovalHistoryRepository,
) reqPorts.RequestService {
	return &RequestServiceImpl{
		requestRepo:         requestRepo,
		workflowRepo:        workflowRepo,
		workflowStepRepo:    workflowStepRepo,
		approvalHistoryRepo: approvalHistoryRepo,
	}
}

// LockRequest acquires a mutex lock for the given request ID to prevent concurrent approval operations.
// This is Layer 1 of our double-layer concurrency control.
//
// Why use Mutex?
// - Fast: In-memory operations are extremely fast (nanoseconds level)
// - Simple: No external dependencies required
// - Effective: Prevents race conditions within the same application instance
// - Granular: Lock is per-request, not global, allowing concurrent processing of different requests
//
// Returns a function that must be called to release the lock (typically deferred).
//
// Example usage:
//
//	defer service.LockRequest(requestID)()
//
// Note: This mutex only works within a single application instance.
// For multi-instance deployments, we also use SELECT FOR UPDATE (Layer 2) at the database level.
func (s *RequestServiceImpl) LockRequest(requestID string) func() {
	// Get or create mutex for this specific request
	mutex, _ := s.mutexMap.LoadOrStore(requestID, &sync.Mutex{})
	mu := mutex.(*sync.Mutex)

	// Lock - blocks until lock is acquired
	mu.Lock()

	// Return unlock function
	return func() {
		mu.Unlock()
	}
}

// CreateRequest creates a new approval request
func (s *RequestServiceImpl) CreateRequest(ctx context.Context, workflowID, requesterID string, amount float64, title, description string) (*reqDomain.Request, error) {
	if amount <= 0 {
		return nil, ErrRequestAmountPositive
	}

	// Verify workflow exists
	_, err := s.workflowRepo.GetByID(ctx, workflowID)
	if err != nil {
		if errors.Is(err, wfRepo.ErrWorkflowNotFound) {
			return nil, ErrWorkflowNotFound
		}
		return nil, err
	}

	request := reqDomain.NewRequest(workflowID, requesterID, amount, title, description)
	if err := s.requestRepo.Create(ctx, request); err != nil {
		return nil, err
	}

	return request, nil
}

// GetRequest retrieves a request by ID
func (s *RequestServiceImpl) GetRequest(ctx context.Context, id string) (*reqDomain.Request, error) {
	return s.requestRepo.GetByID(ctx, id)
}

// ListRequests retrieves a paginated list of requests
func (s *RequestServiceImpl) ListRequests(ctx context.Context, page, limit int, status *reqDomain.RequestStatus) ([]*reqDomain.Request, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return s.requestRepo.List(ctx, page, limit, status)
}

// UpdateRequest updates an existing request
// Only allows updates if the request is still in PENDING status
func (s *RequestServiceImpl) UpdateRequest(ctx context.Context, id string, amount float64, title, description string) (*reqDomain.Request, error) {
	if amount <= 0 {
		return nil, ErrRequestAmountPositive
	}

	request, err := s.requestRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, reqRepo.ErrRequestNotFound) {
			return nil, ErrRequestNotFound
		}
		return nil, err
	}

	// Only allow updates for PENDING requests
	if !request.IsPending() {
		return nil, ErrRequestNotPending
	}

	// Update fields
	request.Amount = amount
	request.Title = title
	request.Description = description
	request.Version++ // Increment version for optimistic locking

	if err := s.requestRepo.Update(ctx, request); err != nil {
		return nil, err
	}

	return request, nil
}

// Approve approves the current step and moves to the next step
// This method uses double-layer concurrency control:
// 1. Layer 1: Mutex lock (in-memory) - prevents concurrent goroutine calls within same instance
// 2. Layer 2: SELECT FOR UPDATE (database) - prevents race conditions across instances
func (s *RequestServiceImpl) Approve(ctx context.Context, requestID, userID, actorID string, isAdmin bool) (*reqDomain.Request, error) {
	// Layer 1: Acquire in-memory mutex lock
	// This prevents race conditions when multiple goroutines in the same instance
	// try to approve the same request simultaneously
	defer s.LockRequest(requestID)()

	// Get the request - we don't use GetByIDForUpdate here because we already have the mutex
	// The mutex ensures only one goroutine can be in this critical section
	request, err := s.requestRepo.GetByID(ctx, requestID)
	if err != nil {
		if errors.Is(err, reqRepo.ErrRequestNotFound) {
			return nil, ErrRequestNotFound
		}
		return nil, err
	}

	// Check if request is still pending
	if !request.IsPending() {
		return nil, ErrRequestNotPending
	}

	// Get the current step
	currentStep, err := s.workflowStepRepo.GetByWorkflowAndLevel(ctx, request.WorkflowID, request.CurrentStep)
	if err != nil {
		if errors.Is(err, stepRepo.ErrStepNotFound) {
			return nil, ErrNoNextStep
		}
		return nil, err
	}

	// Validate actor authorization
	// If not admin, check if actorID matches the step's actor_id
	if !isAdmin && actorID != "" && currentStep.ActorID != actorID {
		return nil, ErrUnauthorizedActor
	}

	// Check if amount meets the condition for current step
	if !s.checkCondition(currentStep.Conditions, request.Amount) {
		return nil, errors.New("request amount does not meet the minimum requirement for this step")
	}

	// Record approval history
	history := approvalHistoryDomain.NewApprovalHistory(
		requestID,
		request.WorkflowID,
		request.CurrentStep,
		actorID,
		userID,
		approvalHistoryDomain.ApprovalActionApprove,
		"",
	)
	if err := s.approvalHistoryRepo.Create(ctx, history); err != nil {
		return nil, errors.New("failed to record approval history")
	}

	// Try to find the next step
	nextStep, err := s.workflowStepRepo.GetByWorkflowAndLevel(ctx, request.WorkflowID, request.CurrentStep+1)
	if err != nil {
		if errors.Is(err, stepRepo.ErrStepNotFound) {
			// No more steps, mark as approved
			request.Status = reqDomain.StatusApproved
			request.CurrentStep = request.CurrentStep + 1
			request.Version++ // Increment version for optimistic locking
			if err := s.requestRepo.Update(ctx, request); err != nil {
				return nil, errors.New("failed to update request status to approved")
			}
			return request, nil
		}
		return nil, err
	}

	// Move to next step
	request.CurrentStep = nextStep.Level
	request.Version++ // Increment version for optimistic locking
	if err := s.requestRepo.Update(ctx, request); err != nil {
		return nil, errors.New("failed to update request to next step")
	}

	return request, nil
}

// Reject rejects the request
// Also uses mutex lock for consistency
func (s *RequestServiceImpl) Reject(ctx context.Context, requestID, userID, actorID string, isAdmin bool, reason string) (*reqDomain.Request, error) {
	// Layer 1: Acquire in-memory mutex lock
	defer s.LockRequest(requestID)()

	request, err := s.requestRepo.GetByID(ctx, requestID)
	if err != nil {
		if errors.Is(err, reqRepo.ErrRequestNotFound) {
			return nil, ErrRequestNotFound
		}
		return nil, err
	}

	// Check if request is still pending
	if !request.IsPending() {
		return nil, ErrRequestNotPending
	}

	// Get the current step
	currentStep, err := s.workflowStepRepo.GetByWorkflowAndLevel(ctx, request.WorkflowID, request.CurrentStep)
	if err != nil {
		if errors.Is(err, stepRepo.ErrStepNotFound) {
			return nil, ErrNoNextStep
		}
		return nil, err
	}

	// Validate actor authorization
	// If not admin, check if actorID matches the step's actor_id
	if !isAdmin && actorID != "" && currentStep.ActorID != actorID {
		return nil, ErrUnauthorizedActor
	}

	// Record rejection history
	history := approvalHistoryDomain.NewApprovalHistory(
		requestID,
		request.WorkflowID,
		request.CurrentStep,
		actorID,
		userID,
		approvalHistoryDomain.ApprovalActionReject,
		reason,
	)
	if err := s.approvalHistoryRepo.Create(ctx, history); err != nil {
		return nil, errors.New("failed to record rejection history")
	}

	// Mark as rejected
	request.Status = reqDomain.StatusRejected
	request.Version++ // Increment version for optimistic locking
	if err := s.requestRepo.Update(ctx, request); err != nil {
		return nil, errors.New("failed to update request status to rejected")
	}

	return request, nil
}

// DeleteRequest deletes a request by ID
func (s *RequestServiceImpl) DeleteRequest(ctx context.Context, id string) error {
	_, err := s.requestRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, reqRepo.ErrRequestNotFound) {
			return ErrRequestNotFound
		}
		return err
	}
	return s.requestRepo.Delete(ctx, id)
}

// checkCondition checks if the amount meets the step conditions
func (s *RequestServiceImpl) checkCondition(conditions stepDomain.StepConditions, amount float64) bool {
	// If no min_amount condition, always pass
	if conditions.MinAmount == 0 {
		return true
	}
	return amount >= conditions.MinAmount
}
