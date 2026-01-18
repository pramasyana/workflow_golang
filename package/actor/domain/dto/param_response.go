package dto

import "workflow-approval/package/actor/domain"

// ActorResponse represents the actor response
type ActorResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Code      string `json:"code"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ToActorResponse converts an Actor to ActorResponse
func ToActorResponse(a *domain.Actor) *ActorResponse {
	if a == nil {
		return nil
	}
	return &ActorResponse{
		ID:        a.ID,
		Name:      a.Name,
		Code:      a.Code,
		CreatedAt: a.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: a.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// ToActorResponseList converts a list of Actors to ActorResponse list
func ToActorResponseList(actors []*domain.Actor) []*ActorResponse {
	result := make([]*ActorResponse, len(actors))
	for i, a := range actors {
		result[i] = ToActorResponse(a)
	}
	return result
}
