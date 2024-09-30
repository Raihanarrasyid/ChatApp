package auth

type AuthControllerImpl struct {
}

func NewAuthControllerImpl() AuthController {
	return &AuthControllerImpl{}
}