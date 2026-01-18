package dto

// CreateWorkflowRequest represents the create workflow request body
type CreateWorkflowRequest struct {
	Name string `json:"name"`
}

// UpdateWorkflowRequest represents the update workflow request body
type UpdateWorkflowRequest struct {
	Name string `json:"name"`
}
