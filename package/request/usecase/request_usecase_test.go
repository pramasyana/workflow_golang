package usecase

import (
	"context"
	"testing"

	approvalHistoryDomain "workflow-approval/package/approval_history/domain"
	approvalHistoryPorts "workflow-approval/package/approval_history/ports"
	reqDomain "workflow-approval/package/request/domain"
	reqPorts "workflow-approval/package/request/ports"
	reqRepo "workflow-approval/package/request/repository"
	wfDomain "workflow-approval/package/workflow/domain"
	wfPorts "workflow-approval/package/workflow/ports"
	wfRepo "workflow-approval/package/workflow/repository"
	stepDomain "workflow-approval/package/workflow_step/domain"
	stepPorts "workflow-approval/package/workflow_step/ports"
	stepRepo "workflow-approval/package/workflow_step/repository"
)

// MockRequestRepository implements RequestRepository for testing
type MockRequestRepository struct {
	requests map[string]*reqDomain.Request
}

func NewMockRequestRepository() *MockRequestRepository {
	return &MockRequestRepository{
		requests: make(map[string]*reqDomain.Request),
	}
}

func (m *MockRequestRepository) Create(ctx context.Context, request *reqDomain.Request) error {
	m.requests[request.ID] = request
	return nil
}

func (m *MockRequestRepository) GetByID(ctx context.Context, id string) (*reqDomain.Request, error) {
	if r, ok := m.requests[id]; ok {
		return r, nil
	}
	return nil, reqRepo.ErrRequestNotFound
}

func (m *MockRequestRepository) GetByIDForUpdate(ctx context.Context, id string) (*reqDomain.Request, error) {
	return m.GetByID(ctx, id)
}

func (m *MockRequestRepository) Update(ctx context.Context, request *reqDomain.Request) error {
	if _, ok := m.requests[request.ID]; ok {
		m.requests[request.ID] = request
		return nil
	}
	return reqRepo.ErrRequestNotFound
}

func (m *MockRequestRepository) Delete(ctx context.Context, id string) error {
	delete(m.requests, id)
	return nil
}

func (m *MockRequestRepository) List(ctx context.Context, page, limit int, status *reqDomain.RequestStatus) ([]*reqDomain.Request, int64, error) {
	return nil, 0, nil
}

// MockWorkflowRepository implements WorkflowRepository for testing
type MockWorkflowRepository struct {
	workflows map[string]*wfDomain.Workflow
}

func NewMockWorkflowRepository() *MockWorkflowRepository {
	return &MockWorkflowRepository{
		workflows: make(map[string]*wfDomain.Workflow),
	}
}

func (m *MockWorkflowRepository) Create(ctx context.Context, workflow *wfDomain.Workflow) error {
	m.workflows[workflow.ID] = workflow
	return nil
}

func (m *MockWorkflowRepository) GetByID(ctx context.Context, id string) (*wfDomain.Workflow, error) {
	if w, ok := m.workflows[id]; ok {
		return w, nil
	}
	return nil, wfRepo.ErrWorkflowNotFound
}

func (m *MockWorkflowRepository) Update(ctx context.Context, workflow *wfDomain.Workflow) error {
	return nil
}

func (m *MockWorkflowRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (m *MockWorkflowRepository) List(ctx context.Context, page, limit int) ([]*wfDomain.Workflow, int64, error) {
	return nil, 0, nil
}

// MockWorkflowStepRepository implements WorkflowStepRepository for testing
type MockWorkflowStepRepository struct {
	steps map[string]map[int]*stepDomain.WorkflowStep
}

func NewMockWorkflowStepRepository() *MockWorkflowStepRepository {
	return &MockWorkflowStepRepository{
		steps: make(map[string]map[int]*stepDomain.WorkflowStep),
	}
}

func (m *MockWorkflowStepRepository) Create(ctx context.Context, step *stepDomain.WorkflowStep) error {
	if m.steps[step.WorkflowID] == nil {
		m.steps[step.WorkflowID] = make(map[int]*stepDomain.WorkflowStep)
	}
	m.steps[step.WorkflowID][step.Level] = step
	return nil
}

func (m *MockWorkflowStepRepository) GetByID(ctx context.Context, id string) (*stepDomain.WorkflowStep, error) {
	return nil, nil
}

func (m *MockWorkflowStepRepository) GetByWorkflowAndLevel(ctx context.Context, workflowID string, level int) (*stepDomain.WorkflowStep, error) {
	if steps, ok := m.steps[workflowID]; ok {
		if step, ok := steps[level]; ok {
			return step, nil
		}
	}
	return nil, stepRepo.ErrStepNotFound
}

func (m *MockWorkflowStepRepository) Update(ctx context.Context, step *stepDomain.WorkflowStep) error {
	return nil
}

func (m *MockWorkflowStepRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (m *MockWorkflowStepRepository) GetAllByWorkflowID(ctx context.Context, workflowID string) ([]*stepDomain.WorkflowStep, error) {
	return nil, nil
}

func (m *MockWorkflowStepRepository) GetByWorkflowID(ctx context.Context, workflowID string) ([]*stepDomain.WorkflowStep, error) {
	if steps, ok := m.steps[workflowID]; ok {
		result := make([]*stepDomain.WorkflowStep, 0, len(steps))
		for _, step := range steps {
			result = append(result, step)
		}
		return result, nil
	}
	return nil, nil
}

// MockApprovalHistoryRepository implements ApprovalHistoryRepository for testing
type MockApprovalHistoryRepository struct {
	histories map[string][]*approvalHistoryDomain.ApprovalHistory
}

func NewMockApprovalHistoryRepository() *MockApprovalHistoryRepository {
	return &MockApprovalHistoryRepository{
		histories: make(map[string][]*approvalHistoryDomain.ApprovalHistory),
	}
}

func (m *MockApprovalHistoryRepository) Create(ctx context.Context, history *approvalHistoryDomain.ApprovalHistory) error {
	m.histories[history.RequestID] = append(m.histories[history.RequestID], history)
	return nil
}

func (m *MockApprovalHistoryRepository) GetByID(ctx context.Context, id string) (*approvalHistoryDomain.ApprovalHistory, error) {
	return nil, nil
}

func (m *MockApprovalHistoryRepository) GetByRequestID(ctx context.Context, requestID string) ([]*approvalHistoryDomain.ApprovalHistory, error) {
	if h, ok := m.histories[requestID]; ok {
		return h, nil
	}
	return nil, nil
}

func (m *MockApprovalHistoryRepository) GetByActorID(ctx context.Context, actorID string) ([]*approvalHistoryDomain.ApprovalHistory, error) {
	return nil, nil
}

func (m *MockApprovalHistoryRepository) GetByRequestIDOrdered(ctx context.Context, requestID string) ([]*approvalHistoryDomain.ApprovalHistory, error) {
	return m.GetByRequestID(ctx, requestID)
}

// Test helper functions
func createTestWorkflow(id string) *wfDomain.Workflow {
	return &wfDomain.Workflow{
		ID:   id,
		Name: "Test Workflow",
	}
}

func createTestStep(workflowID string, level int, minAmount float64, actorID string) *stepDomain.WorkflowStep {
	return &stepDomain.WorkflowStep{
		WorkflowID: workflowID,
		Level:      level,
		ActorID:    actorID,
		Conditions: stepDomain.StepConditions{
			MinAmount: minAmount,
		},
	}
}

func createTestRequest(id, workflowID string, amount float64, step int, status reqDomain.RequestStatus) *reqDomain.Request {
	return &reqDomain.Request{
		ID:          id,
		WorkflowID:  workflowID,
		Amount:      amount,
		CurrentStep: step,
		Status:      status,
		Title:       "Test Request",
	}
}

// Test cases
func TestRequestCreation(t *testing.T) {
	ctx := context.Background()
	mockRequestRepo := NewMockRequestRepository()
	mockWorkflowRepo := NewMockWorkflowRepository()
	mockStepRepo := NewMockWorkflowStepRepository()
	mockApprovalHistoryRepo := NewMockApprovalHistoryRepository()

	// Create workflow first
	workflow := createTestWorkflow("wf-1")
	mockWorkflowRepo.Create(ctx, workflow)

	service := NewRequestService(mockRequestRepo, mockWorkflowRepo, mockStepRepo, mockApprovalHistoryRepo)

	t.Run("Create valid request", func(t *testing.T) {
		req, err := service.CreateRequest(ctx, "wf-1", "user-1", 1500000, "Test Request", "Description")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if req == nil {
			t.Error("Expected request to be created")
		}
		if req.Status != reqDomain.StatusPending {
			t.Errorf("Expected status PENDING, got %s", req.Status)
		}
		if req.CurrentStep != 1 {
			t.Errorf("Expected current step 1, got %d", req.CurrentStep)
		}
	})

	t.Run("Create request with invalid amount", func(t *testing.T) {
		_, err := service.CreateRequest(ctx, "wf-1", "user-1", 0, "Test Request", "Description")
		if err != ErrRequestAmountPositive {
			t.Errorf("Expected ErrRequestAmountPositive, got %v", err)
		}
	})

	t.Run("Create request with negative amount", func(t *testing.T) {
		_, err := service.CreateRequest(ctx, "wf-1", "user-1", -100, "Test Request", "Description")
		if err != ErrRequestAmountPositive {
			t.Errorf("Expected ErrRequestAmountPositive, got %v", err)
		}
	})

	t.Run("Create request with non-existent workflow", func(t *testing.T) {
		_, err := service.CreateRequest(ctx, "non-existent", "user-1", 1500000, "Test Request", "Description")
		if err != ErrWorkflowNotFound {
			t.Errorf("Expected ErrWorkflowNotFound, got %v", err)
		}
	})
}

func TestApprovalLogic(t *testing.T) {
	t.Run("Approve request successfully - single step workflow", func(t *testing.T) {
		ctx := context.Background()
		mockRequestRepo := NewMockRequestRepository()
		mockWorkflowRepo := NewMockWorkflowRepository()
		mockStepRepo := NewMockWorkflowStepRepository()
		mockApprovalHistoryRepo := NewMockApprovalHistoryRepository()

		// Setup workflow with only step 1 (single step workflow)
		workflow := createTestWorkflow("wf-1")
		mockWorkflowRepo.Create(ctx, workflow)

		step1 := createTestStep("wf-1", 1, 1000000, "approver-1")
		mockStepRepo.Create(ctx, step1)

		service := NewRequestService(mockRequestRepo, mockWorkflowRepo, mockStepRepo, mockApprovalHistoryRepo)

		// Create request with amount that exceeds step 1 min_amount
		req := createTestRequest("req-1", "wf-1", 2000000, 1, reqDomain.StatusPending)
		mockRequestRepo.Create(ctx, req)

		approved, err := service.Approve(ctx, "req-1", "user-1", "approver-1", false)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if approved.Status != reqDomain.StatusApproved {
			t.Errorf("Expected status APPROVED, got %s", approved.Status)
		}
	})

	t.Run("Approve request - move to next step", func(t *testing.T) {
		ctx := context.Background()
		mockRequestRepo := NewMockRequestRepository()
		mockWorkflowRepo := NewMockWorkflowRepository()
		mockStepRepo := NewMockWorkflowStepRepository()
		mockApprovalHistoryRepo := NewMockApprovalHistoryRepository()

		// Setup workflow with step 1 and step 2
		workflow := createTestWorkflow("wf-1")
		mockWorkflowRepo.Create(ctx, workflow)

		step1 := createTestStep("wf-1", 1, 1000000, "approver-1")
		mockStepRepo.Create(ctx, step1)

		step2 := createTestStep("wf-1", 2, 5000000, "approver-2")
		mockStepRepo.Create(ctx, step2)

		service := NewRequestService(mockRequestRepo, mockWorkflowRepo, mockStepRepo, mockApprovalHistoryRepo)

		// Create request with amount that exceeds step 1 but not step 2
		req := createTestRequest("req-2", "wf-1", 2000000, 1, reqDomain.StatusPending)
		mockRequestRepo.Create(ctx, req)

		approved, err := service.Approve(ctx, "req-2", "user-1", "approver-1", false)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if approved.Status != reqDomain.StatusPending {
			t.Errorf("Expected status PENDING, got %s", approved.Status)
		}
		if approved.CurrentStep != 2 {
			t.Errorf("Expected current step 2, got %d", approved.CurrentStep)
		}
	})

	t.Run("Approve request - amount below min_amount", func(t *testing.T) {
		ctx := context.Background()
		mockRequestRepo := NewMockRequestRepository()
		mockWorkflowRepo := NewMockWorkflowRepository()
		mockStepRepo := NewMockWorkflowStepRepository()
		mockApprovalHistoryRepo := NewMockApprovalHistoryRepository()

		workflow := createTestWorkflow("wf-1")
		mockWorkflowRepo.Create(ctx, workflow)

		step1 := createTestStep("wf-1", 1, 1000000, "approver-1")
		mockStepRepo.Create(ctx, step1)

		service := NewRequestService(mockRequestRepo, mockWorkflowRepo, mockStepRepo, mockApprovalHistoryRepo)

		// Create request with amount below step 1 min_amount
		req := createTestRequest("req-3", "wf-1", 500000, 1, reqDomain.StatusPending)
		mockRequestRepo.Create(ctx, req)

		_, err := service.Approve(ctx, "req-3", "user-1", "approver-1", false)
		if err == nil {
			t.Error("Expected error for amount below min_amount")
		}
	})

	t.Run("Approve already approved request", func(t *testing.T) {
		ctx := context.Background()
		mockRequestRepo := NewMockRequestRepository()
		mockWorkflowRepo := NewMockWorkflowRepository()
		mockStepRepo := NewMockWorkflowStepRepository()
		mockApprovalHistoryRepo := NewMockApprovalHistoryRepository()

		workflow := createTestWorkflow("wf-1")
		mockWorkflowRepo.Create(ctx, workflow)

		service := NewRequestService(mockRequestRepo, mockWorkflowRepo, mockStepRepo, mockApprovalHistoryRepo)

		req := createTestRequest("req-4", "wf-1", 2000000, 2, reqDomain.StatusApproved)
		mockRequestRepo.Create(ctx, req)

		_, err := service.Approve(ctx, "req-4", "user-1", "approver-1", false)
		if err != ErrRequestNotPending {
			t.Errorf("Expected ErrRequestNotPending, got %v", err)
		}
	})

	t.Run("Approve rejected request", func(t *testing.T) {
		ctx := context.Background()
		mockRequestRepo := NewMockRequestRepository()
		mockWorkflowRepo := NewMockWorkflowRepository()
		mockStepRepo := NewMockWorkflowStepRepository()
		mockApprovalHistoryRepo := NewMockApprovalHistoryRepository()

		workflow := createTestWorkflow("wf-1")
		mockWorkflowRepo.Create(ctx, workflow)

		service := NewRequestService(mockRequestRepo, mockWorkflowRepo, mockStepRepo, mockApprovalHistoryRepo)

		req := createTestRequest("req-5", "wf-1", 2000000, 1, reqDomain.StatusRejected)
		mockRequestRepo.Create(ctx, req)

		_, err := service.Approve(ctx, "req-5", "user-1", "approver-1", false)
		if err != ErrRequestNotPending {
			t.Errorf("Expected ErrRequestNotPending, got %v", err)
		}
	})

	t.Run("Approve non-existent request", func(t *testing.T) {
		ctx := context.Background()
		mockRequestRepo := NewMockRequestRepository()
		mockWorkflowRepo := NewMockWorkflowRepository()
		mockStepRepo := NewMockWorkflowStepRepository()
		mockApprovalHistoryRepo := NewMockApprovalHistoryRepository()

		workflow := createTestWorkflow("wf-1")
		mockWorkflowRepo.Create(ctx, workflow)

		service := NewRequestService(mockRequestRepo, mockWorkflowRepo, mockStepRepo, mockApprovalHistoryRepo)

		_, err := service.Approve(ctx, "non-existent", "user-1", "approver-1", false)
		if err != ErrRequestNotFound {
			t.Errorf("Expected ErrRequestNotFound, got %v", err)
		}
	})
}

func TestRejectLogic(t *testing.T) {
	ctx := context.Background()
	mockRequestRepo := NewMockRequestRepository()
	mockWorkflowRepo := NewMockWorkflowRepository()
	mockStepRepo := NewMockWorkflowStepRepository()
	mockApprovalHistoryRepo := NewMockApprovalHistoryRepository()

	workflow := createTestWorkflow("wf-1")
	mockWorkflowRepo.Create(ctx, workflow)

	// Create workflow step for rejection tests
	step1 := createTestStep("wf-1", 1, 1000000, "approver-1")
	mockStepRepo.Create(ctx, step1)

	service := NewRequestService(mockRequestRepo, mockWorkflowRepo, mockStepRepo, mockApprovalHistoryRepo)

	t.Run("Reject pending request", func(t *testing.T) {
		req := createTestRequest("req-1", "wf-1", 1500000, 1, reqDomain.StatusPending)
		mockRequestRepo.Create(ctx, req)

		rejected, err := service.Reject(ctx, "req-1", "user-1", "approver-1", false, "Insufficient documentation")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if rejected.Status != reqDomain.StatusRejected {
			t.Errorf("Expected status REJECTED, got %s", rejected.Status)
		}
	})

	t.Run("Reject already rejected request", func(t *testing.T) {
		req := createTestRequest("req-2", "wf-1", 1500000, 1, reqDomain.StatusRejected)
		mockRequestRepo.Create(ctx, req)

		_, err := service.Reject(ctx, "req-2", "user-1", "approver-1", false, "Another reason")
		if err != ErrRequestNotPending {
			t.Errorf("Expected ErrRequestNotPending, got %v", err)
		}
	})
}

// Interface compliance tests - ensure mocks implement the interfaces
var _ reqPorts.RequestRepository = (*MockRequestRepository)(nil)
var _ wfPorts.WorkflowRepository = (*MockWorkflowRepository)(nil)
var _ stepPorts.WorkflowStepRepository = (*MockWorkflowStepRepository)(nil)
var _ approvalHistoryPorts.ApprovalHistoryRepository = (*MockApprovalHistoryRepository)(nil)
