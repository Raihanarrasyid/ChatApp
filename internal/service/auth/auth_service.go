package auth

import (
	"ChatApp/internal/http/request"
	"ChatApp/internal/http/response"
	"context"
)

type AuthService interface {
	SendOTP(ctx context.Context, req request.CreateUserRequest) error
	VerifyAndCreateUser(ctx context.Context, req request.VerifyOtpRequest) error
	SignIn(ctx context.Context, req request.SignInRequest, jwtSecret string) (response.SignInResponse, error)
}