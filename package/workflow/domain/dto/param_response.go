package dto

import "workflow-approval/package/workflow/domain"

// WorkflowResponse represents the workflow response
type WorkflowResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

// ToWorkflowResponse converts a Workflow to WorkflowResponse
func ToWorkflowResponse(w *domain.Workflow) *WorkflowResponse {
	if w == nil {
		return nil
	}
	return &WorkflowResponse{
		ID:        w.ID,
		Name:      w.Name,
		CreatedAt: w.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// ToWorkflowResponseList converts a list of Workflow to WorkflowResponse
func ToWorkflowResponseList(workflows []*domain.Workflow) []*WorkflowResponse {
	responses := make([]*WorkflowResponse, len(workflows))
	for i, w := range workflows {
		responses[i] = ToWorkflowResponse(w)
	}
	return responses
}
