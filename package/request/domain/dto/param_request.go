package dto

// CreateRequestRequest represents the create request request body
type CreateRequestRequest struct {
	WorkflowID  string  `json:"workflow_id"`
	Amount      float64 `json:"amount"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
}

// UpdateRequestRequest represents the update request request body
type UpdateRequestRequest struct {
	Amount      float64 `json:"amount"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
}

// RejectRequest represents the reject request body
type RejectRequest struct {
	Reason string `json:"reason"`
}
