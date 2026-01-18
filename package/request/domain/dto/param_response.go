package dto

import (
	"workflow-approval/package/request/domain"
)

// RequestResponse represents the request response
type RequestResponse struct {
	ID          string               `json:"id"`
	WorkflowID  string               `json:"workflow_id"`
	CurrentStep int                  `json:"current_step"`
	Status      domain.RequestStatus `json:"status"`
	Amount      float64              `json:"amount"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	RequesterID string               `json:"requester_id"`
	CreatedAt   string               `json:"created_at"`
	UpdatedAt   string               `json:"updated_at"`
}

// ToRequestResponse converts a Request to RequestResponse
func ToRequestResponse(r *domain.Request) *RequestResponse {
	if r == nil {
		return nil
	}
	return &RequestResponse{
		ID:          r.ID,
		WorkflowID:  r.WorkflowID,
		CurrentStep: r.CurrentStep,
		Status:      r.Status,
		Amount:      r.Amount,
		Title:       r.Title,
		Description: r.Description,
		RequesterID: r.RequesterID,
		CreatedAt:   r.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   r.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// ToRequestResponseList converts a list of Request to RequestResponse
func ToRequestResponseList(requests []*domain.Request) []*RequestResponse {
	responses := make([]*RequestResponse, len(requests))
	for i, r := range requests {
		responses[i] = ToRequestResponse(r)
	}
	return responses
}
