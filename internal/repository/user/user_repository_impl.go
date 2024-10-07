package user

import (
	"ChatApp/internal/model"
	"context"

	commonErr "ChatApp/internal/util/error"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{db}
}

func (ur *UserRepositoryImpl) Create(ctx context.Context, user *model.User) (err error) {
	if err := ur.db.Create(user).Error; err != nil {
		return commonErr.NewBadRequest("Failed to create user: " + err.Error())
	}
	return nil
}

func (ur *UserRepositoryImpl) Update(ctx context.Context, user *model.User) (err error) {
	if err := ur.db.Save(user).Error; err != nil {
		return commonErr.NewBadRequest("Failed to update user: " + err.Error())
	}
	return nil
}

func (ur *UserRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) (err error) {
	if err := ur.db.Delete(&model.User{}, id).Error; err != nil {
		return commonErr.NewBadRequest("Failed to delete user: " + err.Error())
	}
	return nil
}

func (ur *UserRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (user *model.User, err error) {
	if err := ur.db.First(user, id).Error; err != nil {
		return nil, commonErr.NewNotFound("User not found")
	}
	return user, nil
}

func (ur *UserRepositoryImpl) GetAll(ctx context.Context, limit, offset int) (users []*model.User, err error) {
	err = ur.db.Limit(limit).Offset(offset).Find(&users).Error
	if err != nil {
		userNil := make([]*model.User, 0)
		return userNil, nil
	}
	return users, nil
}

func (ur *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (user *model.User, err error) {
	if err := ur.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, commonErr.NewNotFound("User not found")
	}
	return user, nil
}