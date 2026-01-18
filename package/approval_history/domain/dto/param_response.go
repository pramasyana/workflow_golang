package dto

import (
	"workflow-approval/package/approval_history/domain"
)

// ApprovalHistoryResponse represents the approval history response with detailed information
type ApprovalHistoryResponse struct {
	ID           string `json:"id"`
	RequestID    string `json:"request_id"`
	WorkflowID   string `json:"workflow_id"`
	WorkflowName string `json:"workflow_name,omitempty"`
	StepLevel    int    `json:"step_level"`
	StepDesc     string `json:"step_description,omitempty"`
	ActorID      string `json:"actor_id"`
	ActorName    string `json:"actor_name,omitempty"`
	ActorCode    string `json:"actor_code,omitempty"`
	UserID       string `json:"user_id"`
	UserName     string `json:"user_name,omitempty"`
	Action       string `json:"action"`
	Comment      string `json:"comment,omitempty"`
	CreatedAt    string `json:"created_at"`
}

// ToApprovalHistoryResponse converts ApprovalHistory to ApprovalHistoryResponse
func ToApprovalHistoryResponse(h *domain.ApprovalHistory, actorName, actorCode, workflowName, stepDesc string) *ApprovalHistoryResponse {
	return &ApprovalHistoryResponse{
		ID:           h.ID,
		RequestID:    h.RequestID,
		WorkflowID:   h.WorkflowID,
		WorkflowName: workflowName,
		StepLevel:    h.StepLevel,
		StepDesc:     stepDesc,
		ActorID:      h.ActorID,
		ActorName:    actorName,
		ActorCode:    actorCode,
		Action:       string(h.Action),
		Comment:      h.Comment,
		CreatedAt:    h.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

// ToApprovalHistoryResponseList converts a list of ApprovalHistory to ApprovalHistoryResponse
func ToApprovalHistoryResponseList(
	histories []*domain.ApprovalHistory,
	getActorName func(actorID string) (string, string),
	getWorkflowName func(workflowID string) string,
	getStepDesc func(workflowID string, level int) string,
) []*ApprovalHistoryResponse {
	responses := make([]*ApprovalHistoryResponse, len(histories))
	for i, h := range histories {
		actorName, actorCode := getActorName(h.ActorID)
		workflowName := getWorkflowName(h.WorkflowID)
		stepDesc := getStepDesc(h.WorkflowID, h.StepLevel)
		responses[i] = ToApprovalHistoryResponse(h, actorName, actorCode, workflowName, stepDesc)
	}
	return responses
}

// ApprovalHistoryDetailResponse represents detailed approval history for a request
type ApprovalHistoryDetailResponse struct {
	RequestID    string                     `json:"request_id"`
	RequestTitle string                     `json:"request_title"`
	WorkflowID   string                     `json:"workflow_id"`
	WorkflowName string                     `json:"workflow_name"`
	TotalSteps   int                        `json:"total_steps"`
	CurrentStep  int                        `json:"current_step"`
	Status       string                     `json:"status"`
	History      []*ApprovalHistoryResponse `json:"history"`
}
