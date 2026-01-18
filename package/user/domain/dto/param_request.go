package dto

// CreateUserRequest represents the create user request body
type CreateUserRequest struct {
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Name     string  `json:"name"`
	IsAdmin  bool    `json:"is_admin"`
	ActorID  *string `json:"actor_id"`
}

// UpdateProfileRequest represents the update profile request body
type UpdateProfileRequest struct {
	Name string `json:"name"`
}
