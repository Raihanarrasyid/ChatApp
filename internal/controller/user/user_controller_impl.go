package user

type UserControllerImpl struct {
}

func NewUserControllerImpl() UserController {
	return &UserControllerImpl{}
}