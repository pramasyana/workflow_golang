package domain

import (
	"time"

	"workflow-approval/utils"
)

// User represents a user in the system
type User struct {
	ID        string    `json:"id" gorm:"primaryKey;size:36"`
	Email     string    `json:"email" gorm:"uniqueIndex;size:255;not null"`
	Password  string    `json:"-" gorm:"size:255;not null"`
	Name      string    `json:"name" gorm:"size:255;not null"`
	IsAdmin   bool      `json:"is_admin" gorm:"default:false"`
	ActorID   *string   `json:"actor_id" gorm:"size:36"` // Optional, only for non-admin users
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewUser creates a new User instance
func NewUser(email, password, name string, isAdmin bool, actorID *string) *User {
	return &User{
		ID:        utils.GenerateUUID(),
		Email:     email,
		Password:  password,
		Name:      name,
		IsAdmin:   isAdmin,
		ActorID:   actorID,
		CreatedAt: utils.TimeNowUTC(),
		UpdatedAt: utils.TimeNowUTC(),
	}
}

// TableName returns the table name for GORM
func (User) TableName() string {
	return "users"
}
