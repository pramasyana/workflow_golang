package domain

import (
	"time"

	"workflow-approval/utils"
)

// Workflow represents an approval workflow
type Workflow struct {
	ID        string    `json:"id" gorm:"primaryKey;size:36"`
	Name      string    `json:"name" gorm:"size:255;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewWorkflow creates a new Workflow instance
func NewWorkflow(name string) *Workflow {
	now := utils.TimeNowUTC()
	return &Workflow{
		ID:        utils.GenerateUUID(),
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// TableName returns the table name for GORM
func (Workflow) TableName() string {
	return "workflows"
}
