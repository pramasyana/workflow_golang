package domain

import (
	"time"

	"workflow-approval/utils"
)

// RequestStatus represents the status of a workflow request
type RequestStatus string

const (
	StatusPending  RequestStatus = "PENDING"
	StatusApproved RequestStatus = "APPROVED"
	StatusRejected RequestStatus = "REJECTED"
)

// Request represents a workflow approval request
type Request struct {
	ID          string        `json:"id" gorm:"primaryKey;size:36"`
	WorkflowID  string        `json:"workflow_id" gorm:"size:36;not null;index"`
	CurrentStep int           `json:"current_step" gorm:"not null;default:1"`
	Status      RequestStatus `json:"status" gorm:"size:20;not null;default:'PENDING'"`
	Amount      float64       `json:"amount" gorm:"type:decimal(15,2);not null"`
	Title       string        `json:"title" gorm:"size:255"`
	Description string        `json:"description" gorm:"type:text"`
	RequesterID string        `json:"requester_id" gorm:"size:36;not null;index"`
	Version     int           `json:"version" gorm:"not null;default:1"` // For optimistic locking
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

// NewRequest creates a new Request instance
func NewRequest(workflowID, requesterID string, amount float64, title, description string) *Request {
	now := utils.TimeNowUTC()
	return &Request{
		ID:          utils.GenerateUUID(),
		WorkflowID:  workflowID,
		CurrentStep: 1,
		Status:      StatusPending,
		Amount:      amount,
		Title:       title,
		Description: description,
		RequesterID: requesterID,
		Version:     1,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// TableName returns the table name for GORM
func (Request) TableName() string {
	return "requests"
}

// IsPending checks if the request is still pending
func (r *Request) IsPending() bool {
	return r.Status == StatusPending
}

// IsTerminal checks if the request has reached a terminal state
func (r *Request) IsTerminal() bool {
	return r.Status == StatusApproved || r.Status == StatusRejected
}
