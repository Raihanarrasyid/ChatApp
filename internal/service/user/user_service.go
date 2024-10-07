package user

import (
	"ChatApp/internal/http/request"
	"ChatApp/internal/http/response"
	"context"
)

type UserService interface {
	Update(ctx context.Context, id string, req request.UpdateUserRequest) (err error)
	Delete(ctx context.Context, id string) (err error)
	GetByID(ctx context.Context, id string) (user response.UserResponse, err error)
	GetAll(ctx context.Context, page, pageSize int) (users []response.UserResponse, err error)
}