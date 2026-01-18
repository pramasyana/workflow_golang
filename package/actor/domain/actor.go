package domain

import (
	"time"

	"workflow-approval/utils"
)

// Actor represents an actor in the workflow system
type Actor struct {
	ID        string    `json:"id" gorm:"primaryKey;size:36"`
	Name      string    `json:"name" gorm:"size:255;not null"`
	Code      string    `json:"code" gorm:"size:50;uniqueIndex;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewActor creates a new Actor instance
func NewActor(name, code string) *Actor {
	return &Actor{
		ID:        utils.GenerateUUID(),
		Name:      name,
		Code:      code,
		CreatedAt: utils.TimeNowUTC(),
		UpdatedAt: utils.TimeNowUTC(),
	}
}

// TableName returns the table name for GORM
func (Actor) TableName() string {
	return "actors"
}
