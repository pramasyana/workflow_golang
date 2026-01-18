package dto

// CreateActorRequest represents the create actor request body
type CreateActorRequest struct {
	Name string `json:"name" validate:"required"`
	Code string `json:"" validate:"required"`
}

// UpdateActorRequest represents the update actor request body
type UpdateActorRequest struct {
	Name string `json:"name" validate:"required"`
}
