package dto

import "workflow-approval/package/workflow_step/domain"

// CreateStepRequest represents the create step request body
type CreateStepRequest struct {
	Level       int                   `json:"level"`
	ActorID     string                `json:"actor_id"`
	Conditions  domain.StepConditions `json:"conditions"`
	Description string                `json:"description"`
}

// UpdateStepRequest represents the update step request body
type UpdateStepRequest struct {
	Level       int                   `json:"level"`
	ActorID     string                `json:"actor_id"`
	Conditions  domain.StepConditions `json:"conditions"`
	Description string                `json:"description"`
}
