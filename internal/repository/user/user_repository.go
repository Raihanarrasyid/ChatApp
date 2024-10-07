package user

import (
	"ChatApp/internal/model"
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) (err error)
	Update(ctx context.Context, user *model.User) (err error)
	Delete(ctx context.Context, id uuid.UUID) (err error)
	GetByID(ctx context.Context, id uuid.UUID) (user *model.User, err error)
	GetAll(ctx context.Context, limit, offset int) (users []*model.User, err error)
	GetByEmail(ctx context.Context, email string) (user *model.User, err error)
}