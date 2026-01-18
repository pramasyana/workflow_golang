package usecase

import (
	"context"

	"workflow-approval/package/approval_history/domain"
	"workflow-approval/package/approval_history/domain/dto"
	"workflow-approval/package/approval_history/ports"
)

// ApprovalHistoryServiceImpl implements ApprovalHistoryService interface
type ApprovalHistoryServiceImpl struct {
	historyRepo ports.ApprovalHistoryRepository
}

// NewApprovalHistoryService creates a new ApprovalHistoryServiceImpl instance
func NewApprovalHistoryService(historyRepo ports.ApprovalHistoryRepository) ports.ApprovalHistoryService {
	return &ApprovalHistoryServiceImpl{
		historyRepo: historyRepo,
	}
}

// GetHistoryByRequestID retrieves all approval history for a request with details
func (s *ApprovalHistoryServiceImpl) GetHistoryByRequestID(
	ctx context.Context,
	requestID string,
) ([]*dto.ApprovalHistoryResponse, error) {
	histories, err := s.historyRepo.GetByRequestIDOrdered(ctx, requestID)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.ApprovalHistoryResponse, len(histories))
	for i, h := range histories {
		responses[i] = &dto.ApprovalHistoryResponse{
			ID:         h.ID,
			RequestID:  h.RequestID,
			WorkflowID: h.WorkflowID,
			StepLevel:  h.StepLevel,
			ActorID:    h.ActorID,
			UserID:     h.UserID,
			Action:     string(h.Action),
			Comment:    h.Comment,
			CreatedAt:  h.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return responses, nil
}

// CreateApprovalHistory creates a new approval history entry
func (s *ApprovalHistoryServiceImpl) CreateApprovalHistory(
	ctx context.Context,
	requestID, workflowID string,
	stepLevel int,
	actorID, userID string,
	action domain.ApprovalAction,
	comment string,
) error {
	history := domain.NewApprovalHistory(requestID, workflowID, stepLevel, actorID, userID, action, comment)
	return s.historyRepo.Create(ctx, history)
}
