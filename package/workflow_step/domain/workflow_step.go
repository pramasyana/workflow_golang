package domain

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"workflow-approval/utils"
)

// WorkflowStep represents a step in an approval workflow
type WorkflowStep struct {
	ID          string         `json:"id" gorm:"primaryKey;size:36"`
	WorkflowID  string         `json:"workflow_id" gorm:"size:36;not null;index"`
	Level       int            `json:"level" gorm:"not null"`
	ActorID     string         `json:"actor_id" gorm:"size:36;not null"`
	Conditions  StepConditions `json:"conditions" gorm:"type:text"`
	Description string         `json:"description" gorm:"size:500"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// StepConditions represents the conditions for a workflow step
type StepConditions struct {
	MinAmount float64  `json:"min_amount,omitempty"`
	MaxAmount float64  `json:"max_amount,omitempty"`
	Roles     []string `json:"roles,omitempty"`
}

// Value implements driver.Valuer interface for GORM
func (sc StepConditions) Value() (driver.Value, error) {
	// Check if empty by checking if all fields are zero values
	// Use JSON comparison to detect empty struct
	if sc.MinAmount == 0 && sc.MaxAmount == 0 && len(sc.Roles) == 0 {
		// Return empty object instead of nil to preserve field
		return []byte(`{}`), nil
	}
	return json.Marshal(sc)
}

// Scan implements sql.Scanner interface for GORM
func (sc *StepConditions) Scan(value interface{}) error {
	if value == nil {
		*sc = StepConditions{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("type assertion to []byte or string failed")
	}

	return json.Unmarshal(bytes, sc)
}

// NullStepConditions is a wrapper for nullable StepConditions in SQL
type NullStepConditions struct {
	StepConditions
	Valid bool
}

// Value implements driver.Valuer interface
func (ns NullStepConditions) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.StepConditions.Value()
}

// Scan implements sql.Scanner interface
func (ns *NullStepConditions) Scan(value interface{}) error {
	if value == nil {
		ns.StepConditions, ns.Valid = StepConditions{}, false
		return nil
	}

	ns.Valid = true
	return ns.StepConditions.Scan(value)
}

// Ensure StepConditions implements sql.Scanner and driver.Valuer
var _ sql.Scanner = (*StepConditions)(nil)
var _ driver.Valuer = StepConditions{}

// NewWorkflowStep creates a new WorkflowStep instance
func NewWorkflowStep(workflowID string, level int, actorID string, conditions StepConditions) *WorkflowStep {
	now := utils.TimeNowUTC()
	return &WorkflowStep{
		ID:         utils.GenerateUUID(),
		WorkflowID: workflowID,
		Level:      level,
		ActorID:    actorID,
		Conditions: conditions,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// TableName returns the table name for GORM
func (WorkflowStep) TableName() string {
	return "workflow_steps"
}
