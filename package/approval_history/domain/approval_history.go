package domain

import (
	"time"

	"workflow-approval/utils"
)

// ApprovalAction represents the type of approval action
type ApprovalAction string

const (
	ApprovalActionApprove ApprovalAction = "APPROVE"
	ApprovalActionReject  ApprovalAction = "REJECT"
)

// ApprovalHistory represents an approval/rejection history entry
// Tracks who approved/rejected, when, which workflow step, and any comments
type ApprovalHistory struct {
	ID         string         `json:"id" gorm:"primaryKey;size:36"`
	RequestID  string         `json:"request_id" gorm:"size:36;not null;index"`
	WorkflowID string         `json:"workflow_id" gorm:"size:36;not null"`
	StepLevel  int            `json:"step_level" gorm:"not null"`
	ActorID    string         `json:"actor_id" gorm:"size:36;not null;index"`
	UserID     string         `json:"user_id" gorm:"size:36;not null;index"`
	Action     ApprovalAction `json:"action" gorm:"size:20;not null"`
	Comment    string         `json:"comment" gorm:"type:text"`
	CreatedAt  time.Time      `json:"created_at"`
}

// NewApprovalHistory creates a new ApprovalHistory instance
func NewApprovalHistory(requestID, workflowID string, stepLevel int, actorID, userID string, action ApprovalAction, comment string) *ApprovalHistory {
	return &ApprovalHistory{
		ID:         utils.GenerateUUID(),
		RequestID:  requestID,
		WorkflowID: workflowID,
		StepLevel:  stepLevel,
		ActorID:    actorID,
		UserID:     userID,
		Action:     action,
		Comment:    comment,
		CreatedAt:  utils.TimeNowUTC(),
	}
}

// TableName returns the table name for GORM
func (ApprovalHistory) TableName() string {
	return "approval_history"
}
