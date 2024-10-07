package user

import (
	"ChatApp/internal/http/request"
	"ChatApp/internal/http/response"
	"ChatApp/internal/repository/user"
	"context"

	commonError "ChatApp/internal/util/error"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	UserRepository user.UserRepository
}

func NewUserService(userRepository user.UserRepository) UserService {
	return &UserServiceImpl{UserRepository: userRepository}
}

func (usi *UserServiceImpl) Update(ctx context.Context, id string, req request.UpdateUserRequest) (err error) {
	user, err := usi.UserRepository.GetByID(ctx, uuid.MustParse(id))
	if err != nil {
		return commonError.NewNotFound("User not found")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	
	if err != nil {
		return commonError.NewBadRequest("Failed to hash password: " + err.Error())
	}

	user.Username = req.Username
	user.Password = string(hashedPassword)
	user.Email = req.Email

	if err := usi.UserRepository.Update(ctx, user); err != nil {
		return commonError.NewBadRequest("Failed to update user: " + err.Error())
	}

	return nil
}

func (usi *UserServiceImpl) Delete(ctx context.Context, id string) (err error) {
	if err := usi.UserRepository.Delete(ctx, uuid.MustParse(id)); err != nil {
		return err
	}

	return nil
}

func (usi *UserServiceImpl) GetByID(ctx context.Context, id string) (user response.UserResponse, err error) {
	userModel, err := usi.UserRepository.GetByID(ctx, uuid.MustParse(id))
	if err != nil {
		return user, err
	}

	user = response.UserResponse{
		ID:       userModel.ID,
		Username: userModel.Username,
		Email:    userModel.Email,
	}

	return user, nil
}

func (usi *UserServiceImpl) GetAll(ctx context.Context, page, pageSize int) (users []response.UserResponse, err error) {
	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	usersModel, err := usi.UserRepository.GetAll(ctx, pageSize, offset)
	if err != nil {
		return nil, err
	}

	for _, user := range usersModel {
		users = append(users, response.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		})
	}

	return users, nil
}
