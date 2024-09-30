package auth

type AuthRepositoryImpl struct {
}

func NewAuthRepositoryImpl() AuthRepository {
	return &AuthRepositoryImpl{}
}