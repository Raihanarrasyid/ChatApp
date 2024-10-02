package http

type (
	Response struct {
		Message string `json:"message" binding:"required"`
		Value   any    `json:"value"`
	}

	Error struct {
		Message string `json:"message" binding:"required"`
	}
)