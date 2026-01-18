package dto

import "workflow-approval/package/user/domain"

// UserResponse represents the user response
type UserResponse struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	IsAdmin bool   `json:"is_admin"`
}

// ToUserResponse converts a User to UserResponse
func ToUserResponse(u *domain.User) *UserResponse {
	if u == nil {
		return nil
	}
	return &UserResponse{
		ID:      u.ID,
		Email:   u.Email,
		Name:    u.Name,
		IsAdmin: u.IsAdmin,
	}
}
