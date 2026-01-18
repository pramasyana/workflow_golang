package dto

import (
	"workflow-approval/package/workflow_step/domain"
)

// StepResponse represents the step response
type StepResponse struct {
	ID          string                `json:"id"`
	WorkflowID  string                `json:"workflow_id"`
	Level       int                   `json:"level"`
	ActorID     string                `json:"actor_id"`
	Conditions  domain.StepConditions `json:"conditions"`
	Description string                `json:"description"`
	CreatedAt   string                `json:"created_at"`
}

// ToStepResponse converts a WorkflowStep to StepResponse
func ToStepResponse(s *domain.WorkflowStep) *StepResponse {
	if s == nil {
		return nil
	}
	return &StepResponse{
		ID:         s.ID,
		WorkflowID: s.WorkflowID,
		Level:      s.Level,
		ActorID:    s.ActorID,
		Conditions: s.Conditions,
		CreatedAt:  s.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// ToStepResponseList converts a list of WorkflowStep to StepResponse
func ToStepResponseList(steps []*domain.WorkflowStep) []*StepResponse {
	responses := make([]*StepResponse, len(steps))
	for i, s := range steps {
		responses[i] = ToStepResponse(s)
	}
	return responses
}
