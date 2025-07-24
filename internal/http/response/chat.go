package response

import "github.com/google/uuid"

type ChatResponse struct {
	ID uuid.UUID `json:"id"`
	FromID uuid.UUID `json:"from_id"`
	ToID uuid.UUID `json:"to_id"`
	Message string `json:"message"`
	IsRead bool `json:"is_read"`
	FromUser UserResponse `json:"from_user"`
	ToUser UserResponse `json:"to_user"`
}